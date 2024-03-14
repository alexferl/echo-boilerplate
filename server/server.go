package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	casbinMw "github.com/alexferl/echo-casbin"
	jwtMw "github.com/alexferl/echo-jwt"
	openapiMw "github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	jwx "github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	_ "go.uber.org/automaxprocs"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/mappers"
	"github.com/alexferl/echo-boilerplate/services"
	"github.com/alexferl/echo-boilerplate/util/hash"
	"github.com/alexferl/echo-boilerplate/util/jwt"
)

var (
	ErrBanned            = errors.New("account banned")
	ErrLocked            = errors.New("account locked")
	ErrCookieMissing     = errors.New("missing access token cookie")
	ErrCSRFHeaderMissing = errors.New("missing CSRF token header")
	ErrCSRFInvalid       = errors.New("invalid CSRF token")
	ErrTokenInvalid      = errors.New("token invalid")
	ErrTokenMismatch     = errors.New("token mismatch")
	ErrTokenRevoked      = errors.New("token is revoked")
)

func New() *server.Server {
	client, err := data.MewMongoClient()
	if err != nil {
		log.Panic().Err(err).Msg("failed creating mongo client")
	}

	openapi := openapiMw.NewHandler()

	patMapper := mappers.NewPersonalAccessToken(client)
	patSvc := services.NewPersonalAccessToken(patMapper)

	taskMapper := mappers.NewTask(client)
	taskSvc := services.NewTask(taskMapper)

	userMapper := mappers.NewUser(client)
	userSvc := services.NewUser(userMapper)

	return newServer(userSvc, patSvc, []handlers.Handler{
		handlers.NewRootHandler(openapi),
		handlers.NewAuthHandler(openapi, userSvc),
		handlers.NewPersonalAccessTokenHandler(openapi, patSvc),
		handlers.NewTaskHandler(openapi, taskSvc),
		handlers.NewUserHandler(openapi, userSvc),
	}...)
}

func NewTestServer(userSvc handlers.UserService, patSvc handlers.PersonalAccessTokenService, handler ...handlers.Handler) *server.Server {
	c := config.New()
	c.BindFlags()

	viper.Set(config.CookiesEnabled, true)
	viper.Set(config.CSRFEnabled, true)

	return newServer(userSvc, patSvc, handler...)
}

func newServer(userSvc handlers.UserService, patSvc handlers.PersonalAccessTokenService, handler ...handlers.Handler) *server.Server {
	key, err := jwt.LoadPrivateKey()
	if err != nil {
		log.Panic().Err(err).Msg("failed loading private key")
	}

	jwtConfig := jwtMw.Config{
		Key:             key,
		UseRefreshToken: true,
		ExemptRoutes: map[string][]string{
			"/":                       {http.MethodGet},
			"/readyz":                 {http.MethodGet},
			"/livez":                  {http.MethodGet},
			"/docs":                   {http.MethodGet},
			"/openapi/*":              {http.MethodGet},
			"/auth/login":             {http.MethodPost},
			"/auth/signup":            {http.MethodPost},
			"/google":                 {http.MethodGet},
			"/oauth2/google/callback": {http.MethodGet},
			"/oauth2/google/login":    {http.MethodGet},
		},
		AfterParseFunc: func(c echo.Context, t jwx.Token, encodedToken string, src jwtMw.TokenSource) *echo.HTTPError {
			ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
			defer cancel()

			user, err := userSvc.Read(ctx, t.Subject())
			if err != nil {
				log.Error().Err(err).Msg("failed getting user")
				return echo.NewHTTPError(http.StatusServiceUnavailable)
			}

			c.Set("user", user)
			// set roles for casbin
			c.Set("roles", user.Roles)

			if user.IsBanned {
				return echo.NewHTTPError(http.StatusForbidden, ErrBanned.Error())
			}
			if user.IsLocked {
				return echo.NewHTTPError(http.StatusForbidden, ErrLocked.Error())
			}

			// CSRF
			if viper.GetBool(config.CookiesEnabled) && viper.GetBool(config.CSRFEnabled) {
				if src == jwtMw.Cookie {
					switch c.Request().Method {
					case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
					default: // Validate token only for requests which are not defined as 'safe' by RFC7231
						cookie, err := c.Cookie(viper.GetString(config.JWTAccessTokenCookieName))
						if err != nil {
							return echo.NewHTTPError(http.StatusBadRequest, ErrCookieMissing)
						}

						h := c.Request().Header.Get(viper.GetString(config.CSRFHeaderName))
						if h == "" {
							return echo.NewHTTPError(http.StatusBadRequest, ErrCSRFHeaderMissing)
						}

						if !hash.ValidMAC([]byte(cookie.Value), []byte(h), []byte(viper.GetString(config.CSRFSecretKey))) {
							return echo.NewHTTPError(http.StatusForbidden, ErrCSRFInvalid)
						}
					}
				}
			}
			// Personal Access Tokens
			claims := t.PrivateClaims()
			typ := claims["type"]
			if typ == jwt.PersonalToken.String() {
				pat, err := patSvc.FindOne(ctx, t.Subject(), "")
				if err != nil {
					var se *services.Error
					if errors.As(err, &se) {
						if se.Kind == services.NotExist {
							return echo.NewHTTPError(http.StatusUnauthorized, ErrTokenInvalid)
						}
					}
					return echo.NewHTTPError(http.StatusServiceUnavailable)
				}

				if err = pat.Validate(encodedToken); err != nil {
					return echo.NewHTTPError(http.StatusUnauthorized, ErrTokenMismatch)
				}

				if pat.IsRevoked {
					return echo.NewHTTPError(http.StatusUnauthorized, ErrTokenRevoked)
				}
			}

			log.Logger = log.Logger.With().Str("token_id", t.Subject()).Logger()

			return nil
		},
	}

	enforcer, err := casbin.NewEnforcer(viper.GetString(config.CasbinModel), viper.GetString(config.CasbinPolicy))
	if err != nil {
		panic(err)
	}

	openAPIConfig := openapiMw.Config{
		Schema: viper.GetString(config.OpenAPISchema),
		ExemptRoutes: map[string][]string{
			"/":                       {http.MethodGet},
			"/readyz":                 {http.MethodGet},
			"/livez":                  {http.MethodGet},
			"/docs":                   {http.MethodGet},
			"/openapi/*":              {http.MethodGet},
			"/google":                 {http.MethodGet},
			"/oauth2/google/callback": {http.MethodGet},
			"/oauth2/google/login":    {http.MethodGet},
		},
	}

	s := server.New()

	s.Use(
		jwtMw.JWTWithConfig(jwtConfig),
		casbinMw.Casbin(enforcer),
		openapiMw.OpenAPIWithConfig(openAPIConfig),
	)

	for _, h := range handler {
		h.Register(s)
	}

	s.File("/docs", "./docs/index.html")
	s.Static("/openapi/", "./openapi")

	return s
}
