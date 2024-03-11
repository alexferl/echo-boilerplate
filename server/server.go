package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	casbinMw "github.com/alexferl/echo-casbin"
	jwtMw "github.com/alexferl/echo-jwt"
	openapiMw "github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/database/mongodb"
	"github.com/alexferl/golib/http/api/server"
	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	jwx "github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.uber.org/automaxprocs"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/mappers"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/services"
	"github.com/alexferl/echo-boilerplate/util/hash"
	"github.com/alexferl/echo-boilerplate/util/jwt"
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

	return newServer([]handlers.Handler{
		handlers.NewRootHandler(openapi),
		handlers.NewAuthHandler(openapi, userSvc),
		handlers.NewPersonalAccessTokenHandler(openapi, patSvc),
		handlers.NewTaskHandler(openapi, taskSvc),
		handlers.NewUserHandler(openapi, userSvc),
	}...)
}

func NewTestServer(handler ...handlers.Handler) *server.Server {
	c := config.New()
	c.BindFlags()

	viper.Set(config.CookiesEnabled, true)
	viper.Set(config.CSRFEnabled, true)

	return newServer(handler...)
}

func newServer(handler ...handlers.Handler) *server.Server {
	key, err := jwt.LoadPrivateKey()
	if err != nil {
		log.Panic().Err(err).Msg("failed loading private key")
	}

	client, err := mongodb.New()
	if err != nil {
		log.Panic().Err(err).Msg("failed creating mongo client")
	}
	mapper := data.NewMapper(client, viper.GetString(config.AppName), "personal_access_tokens")

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
			// set roles for casbin
			claims := t.PrivateClaims()
			c.Set("roles", claims["roles"])
			isBanned := claims["is_banned"]
			isLocked := claims["is_locked"]

			if val, ok := isBanned.(bool); ok {
				if val {
					return echo.NewHTTPError(http.StatusForbidden, "account banned")
				}
			}
			if val, ok := isLocked.(bool); ok {
				if val {
					return echo.NewHTTPError(http.StatusForbidden, "account locked")
				}
			}

			// CSRF
			if viper.GetBool(config.CookiesEnabled) && viper.GetBool(config.CSRFEnabled) {
				if src == jwtMw.Cookie {
					switch c.Request().Method {
					case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
					default: // Validate token only for requests which are not defined as 'safe' by RFC7231
						cookie, err := c.Cookie(viper.GetString(config.JWTAccessTokenCookieName))
						if err != nil {
							return echo.NewHTTPError(http.StatusBadRequest, "missing access token cookie")
						}

						h := c.Request().Header.Get(viper.GetString(config.CSRFHeaderName))
						if h == "" {
							return echo.NewHTTPError(http.StatusBadRequest, "missing CSRF token header")
						}

						if !hash.ValidMAC([]byte(cookie.Value), []byte(h), []byte(viper.GetString(config.CSRFSecretKey))) {
							return echo.NewHTTPError(http.StatusForbidden, "invalid CSRF token")
						}
					}
				}
			}

			// Personal Access Tokens
			typ := claims["type"]
			if typ == jwt.PersonalToken.String() {
				ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
				defer cancel()

				filter := bson.D{{"user_id", t.Subject()}}
				result, err := mapper.FindOne(ctx, filter, &models.PersonalAccessToken{})
				if err != nil {
					if errors.Is(err, data.ErrNoDocuments) {
						return echo.NewHTTPError(http.StatusUnauthorized, "token invalid")
					}
					return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
				}

				pat := result.(*models.PersonalAccessToken)
				if err = pat.Validate(encodedToken); err != nil {
					return echo.NewHTTPError(http.StatusUnauthorized, "token mismatch")
				}

				if pat.IsRevoked {
					return echo.NewHTTPError(http.StatusUnauthorized, "token is revoked")
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