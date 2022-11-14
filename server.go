package app

import (
	"context"
	"net/http"
	"time"

	casbinMw "github.com/alexferl/echo-casbin"
	jwtMw "github.com/alexferl/echo-jwt"
	openapiMw "github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/handler"
	"github.com/alexferl/golib/http/router"
	"github.com/alexferl/golib/http/server"
	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/handlers/tasks"
	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
)

func DefaultHandlers() []handler.Handler {
	client, err := data.NewClient()
	if err != nil {
		panic(err)
	}
	data.CreateIndexes(client)

	openapi := openapiMw.NewHandler()

	return []handler.Handler{
		handlers.NewHandler(),
		tasks.NewHandler(client, openapi, nil),
		users.NewHandler(client, openapi, nil),
	}
}

func NewServer() *server.Server {
	return newServer(DefaultHandlers()...)
}

func NewTestServer(handler ...handler.Handler) *server.Server {
	c := config.New()
	c.BindFlags()
	return newServer(handler...)
}

func newServer(handler ...handler.Handler) *server.Server {
	var routes []*router.Route
	for _, h := range handler {
		routes = append(routes, h.GetRoutes()...)
	}

	r := &router.Router{Routes: routes}

	key, err := util.LoadPrivateKey()
	if err != nil {
		panic(err)
	}

	client, err := data.NewClient()
	if err != nil {
		panic(err)
	}
	mapper := users.NewMapper(client, users.PATCollection)

	jwtConfig := jwtMw.Config{
		Key:             key,
		UseRefreshToken: true,
		ExemptRoutes: map[string][]string{
			"/":                {http.MethodGet},
			"/healthz":         {http.MethodGet},
			"/favicon.ico":     {http.MethodGet},
			"/docs":            {http.MethodGet},
			"/openapi/*":       {http.MethodGet},
			"/auth/signup":     {http.MethodPost},
			"/auth/login":      {http.MethodPost},
			"/oauth2/login":    {http.MethodGet},
			"/oauth2/callback": {http.MethodGet},
		},
		OptionalRoutes: map[string][]string{
			"/users/:username": {http.MethodGet},
		},
		AfterParseFunc: func(c echo.Context, t jwt.Token) *echo.HTTPError {
			// set roles for casbin
			claims := t.PrivateClaims()
			c.Set("roles", claims["roles"])

			typ := claims["type"]
			if typ == util.PersonalToken.String() {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				filter := bson.D{{"user_id", t.Subject()}}
				result, err := mapper.Collection(users.PATCollection).FindOne(ctx, filter, &users.PersonalAccessToken{})
				if err != nil {
					return echo.NewHTTPError(500, "Internal Server Error")
				}

				pat := result.(*users.PersonalAccessToken)
				if pat.Revoked {
					return echo.NewHTTPError(401, "Token is revoked")
				}
			}

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
			"/":                {http.MethodGet},
			"/healthz":         {http.MethodGet},
			"/favicon.ico":     {http.MethodGet},
			"/docs":            {http.MethodGet},
			"/openapi/*":       {http.MethodGet},
			"/oauth2/login":    {http.MethodGet},
			"/oauth2/callback": {http.MethodGet},
		},
	}

	s := server.New(
		r,
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowCredentials: true,
		}),
		jwtMw.JWTWithConfig(jwtConfig),
		casbinMw.Casbin(enforcer),
		openapiMw.OpenAPIWithConfig(openAPIConfig),
	)

	s.File("/favicon.ico", "./assets/images/favicon.ico")
	s.File("/docs", "./assets/index.html")
	s.Static("/openapi/", "./openapi")

	s.HideBanner = true
	s.HidePort = true

	return s
}
