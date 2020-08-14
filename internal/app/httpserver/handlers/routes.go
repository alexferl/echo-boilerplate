package handlers

import (
	"github.com/labstack/echo/v4"

	"echo-boilerplate/internal/app/httpserver/handlers/healthz"
	"echo-boilerplate/internal/app/httpserver/handlers/root"
)

type (
	// Handler represents the structure of our resource
	Handler struct {
	}
)

// ErrorResponse holds an error message
type ErrorResponse struct {
	Message string `json:"error"`
}

// Register routes with echo
func Register(e *echo.Echo) {
	e.GET("/", root.Root)
	e.GET("/healthz", healthz.Healthz)
}
