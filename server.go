package main

import (
	"github.com/admiralobvious/echo-logrusmiddleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"echo-boilerplate/handlers"
)

func init() {
	cnf := NewConfig()
	cnf.BindFlags()
	InitLogging()
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.POST, echo.OPTIONS, echo.DELETE, echo.PUT},
	}))

	if !viper.GetBool("log-requests-disabled") {
		e.Logger = logrusmiddleware.Logger{Logger: logrus.StandardLogger()}
		e.Use(logrusmiddleware.Hook())
	}

	// Routes
	h := &handlers.Handler{}
	e.GET("/", h.Root)

	// Start server
	addr := viper.GetString("address") + ":" + viper.GetString("port")
	e.Logger.Fatal(e.Start(addr))
}
