package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// Root returns the welcome message
func (h *Handler) Root(c echo.Context) error {
	m := fmt.Sprintf("Welcome to %s", viper.GetString("app-name"))
	return c.JSON(http.StatusOK, Response{m})
}
