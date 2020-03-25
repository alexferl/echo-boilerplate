package main

import (
	"context"
	"os"
	"os/signal"
	"time"

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

	// Routes
	h := &handlers.Handler{}
	e.GET("/", h.Root)

	// Start server
	go func() {
		if err := e.Start(viper.GetString("address") + ":" + viper.GetString("port")); err != nil {
			e.Logger.Info("Received SIGINT, shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
