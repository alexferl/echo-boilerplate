package handlers

import (
	"fmt"
	"net/http"

	libConfig "github.com/alexferl/golib/config"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// Root returns the welcome message.
func (h *Handler) Root(c echo.Context) error {
	m := fmt.Sprintf("Welcome to %s", viper.GetString(libConfig.AppName))
	return c.JSON(http.StatusOK, echo.Map{"message": m})
}
