package app

import (
	"net/http"

	casbinmw "github.com/alexferl/echo-casbin"
	jwtmw "github.com/alexferl/echo-jwt"
	openapimw "github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/handler"
	"github.com/alexferl/golib/http/router"
	"github.com/alexferl/golib/http/server"
	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/spf13/viper"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	hs "github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
)

func DefaultHandlers() []handler.Handler {
	client, err := data.NewClient()
	if err != nil {
		panic(err)
	}

	openapi := openapimw.NewHandler()

	return []handler.Handler{
		hs.NewHandler(),
		users.NewHandler(client, openapi, nil),
	}
}

func NewServer() *server.Server {
	handlers := DefaultHandlers()
	return NewServerWithOverrides(nil, handlers...)
}

func NewServerWithOverrides(overrides map[string]any, handlers ...handler.Handler) *server.Server {
	if overrides != nil {
		for k, v := range overrides {
			viper.Set(k, v)
		}
	}

	var routes []*router.Route
	for _, h := range handlers {
		routes = append(routes, h.GetRoutes()...)
	}

	r := &router.Router{Routes: routes}

	key, err := util.LoadPrivateKey()
	if err != nil {
		panic(err)
	}

	jwtConfig := jwtmw.Config{
		Key:             key,
		UseRefreshToken: true,
		ExemptRoutes: map[string][]string{
			"/":                {http.MethodGet},
			"/healthz":         {http.MethodGet},
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
			return nil
		},
	}

	enforcer, err := casbin.NewEnforcer(viper.GetString(config.CasbinModel), viper.GetString(config.CasbinPolicy))
	if err != nil {
		panic(err)
	}

	openAPIConfig := openapimw.Config{
		Schema: viper.GetString(config.OpenAPISchema),
		ExemptRoutes: map[string][]string{
			"/":        {http.MethodGet},
			"/healthz": {http.MethodGet},
		},
	}

	s := server.New(
		r,
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"http://locahost:1323"},
			AllowCredentials: true,
		}),
		jwtmw.JWTWithConfig(jwtConfig),
		casbinmw.Casbin(enforcer),
		openapimw.OpenAPIWithConfig(openAPIConfig),
	)

	s.HideBanner = true
	s.HidePort = true

	return s
}
