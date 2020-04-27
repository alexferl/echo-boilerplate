package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"github.com/admiralobvious/echo-boilerplate/internal/config"
	"github.com/admiralobvious/echo-boilerplate/internal/handlers"
	"github.com/admiralobvious/echo-boilerplate/internal/logging"
	"github.com/admiralobvious/echo-boilerplate/internal/middleware"
)

func init() {
	cnf := config.NewConfig()
	cnf.BindFlags()
	logging.InitLogging()
}

func main() {
	e := echo.New()

	middleware.Register(e)
	handlers.Register(e)

	// Start server
	go func() {
		address := fmt.Sprintf("%s:%s", viper.GetString("bind-address"), viper.GetString("bind-port"))
		if err := e.Start(address); err != nil {
			e.Logger.Info("Received SIGINT, shutting down the server")
		}
	}()

	timeout := time.Duration(viper.GetInt64("graceful-timeout")) * time.Second

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
