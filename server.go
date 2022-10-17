package app

import (
	"net/http"

	"github.com/alexferl/golib/http/router"
	"github.com/alexferl/golib/http/server"
	"github.com/labstack/echo/v4/middleware"

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
		middleware.BodyLimit("10xM"),
		// more middlewares...
	)

	s.Start()
}
