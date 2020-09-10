package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"echo-boilerplate/internal/app/httpd/router"
)

type Handler struct{}

var handler Handler

// Register routes with echo
func Register(e *echo.Echo) *Handler {
	for _, route := range Router.Routes {
		e.Add(route.Method, route.Pattern, route.HandlerFunc)
	}

	return &handler
}

var Router = &router.Router{
	Routes: []router.Route{
		{
			"Root",
			http.MethodGet,
			"/",
			handler.Root,
		},
		{
			"Healthz",
			http.MethodGet,
			"/healthz",
			handler.Healthz,
		},
	},
}
