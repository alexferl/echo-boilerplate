package root

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// Root returns the welcome message
func Root(c echo.Context) error {
	m := fmt.Sprintf("Welcome to %s", viper.GetString("app-name"))
	log.Info().Msg("derp")
	return c.JSON(http.StatusOK, map[string]string{"message": m})
}
