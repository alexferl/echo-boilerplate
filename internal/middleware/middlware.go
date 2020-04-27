package middleware

import (
	logrusmiddleware "github.com/admiralobvious/echo-logrusmiddleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Register middleware with echo
func Register(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     viper.GetStringSlice("cors-allow-origins"),
		AllowMethods:     viper.GetStringSlice("cors-allow-methods"),
		AllowHeaders:     viper.GetStringSlice("cors-allow-headers"),
		AllowCredentials: viper.GetBool("cors-allow-credentials"),
		ExposeHeaders:    viper.GetStringSlice("cors-expose-headers"),
		MaxAge:           viper.GetInt("cors-max-age"),
	}))

	if !viper.GetBool("log-requests-disabled") {
		e.Logger = logrusmiddleware.Logger{Logger: logrus.StandardLogger()}
		e.Use(logrusmiddleware.Hook())
	}
}
