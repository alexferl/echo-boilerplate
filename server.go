package app

import (
	"net/http"
	"strings"

	"github.com/alexferl/golib/http/router"
	"github.com/alexferl/golib/http/server"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/alexferl/echo-boilerplate/handlers"
)

func Start() {
	c := NewConfig()
	c.BindFlags()

	h := &handlers.Handler{
		// add stuff that the handlers should have access to
		// like a database client.
	}
	r := &router.Router{
		Routes: []*router.Route{
			{"Root", http.MethodGet, "/", h.Root},
			{"Healthz", http.MethodGet, "/healthz", h.Healthz},
		},
	}

	s := server.New(
		r,
		middleware.BodyLimit("1M"),
		// add your own middlewares here
	)

	log.Info().Msgf(
		"Starting %s on %s environment",
		viper.GetString("app-name"),
		strings.ToUpper(viper.GetString("env-name")),
	)

	s.Start()
}
