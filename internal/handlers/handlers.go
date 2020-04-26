package handlers

import (
	"github.com/labstack/echo/v4"
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

func Register(e *echo.Echo) {
	h := &Handler{}
	e.GET("/", h.Root)
}
