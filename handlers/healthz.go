package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Healthz returns the welcome message
func (h *Handler) Healthz(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
