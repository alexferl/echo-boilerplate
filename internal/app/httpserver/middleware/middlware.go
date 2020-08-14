package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"github.com/ziflex/lecho/v2"
)

// Register middleware with echo
func Register(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	if viper.GetBool("cors-enabled") {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     viper.GetStringSlice("cors-allow-origins"),
			AllowMethods:     viper.GetStringSlice("cors-allow-methods"),
			AllowHeaders:     viper.GetStringSlice("cors-allow-headers"),
			AllowCredentials: viper.GetBool("cors-allow-credentials"),
			ExposeHeaders:    viper.GetStringSlice("cors-expose-headers"),
			MaxAge:           viper.GetInt("cors-max-age"),
		}))
	}

	if !viper.GetBool("log-requests-disabled") {
		logger := lecho.New(
			os.Stdout,
			lecho.WithCaller(),
			lecho.WithTimestamp(),
			lecho.WithLevel(log.INFO),
		)
		e.Logger = logger
		e.Use(lecho.Middleware(lecho.Config{
			Logger: logger,
		}))
	}
}
