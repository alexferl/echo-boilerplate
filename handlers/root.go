package handlers

import (
	"fmt"
	"net/http"

	libconfig "github.com/alexferl/golib/config"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type Response struct {
	Message string `json:"message"`
}

// Root returns the welcome message.
func (h *Handler) Root(c echo.Context) error {
	m := fmt.Sprintf("Welcome to %s", viper.GetString(libconfig.AppName))
	return c.JSON(http.StatusOK, Response{m})
}
