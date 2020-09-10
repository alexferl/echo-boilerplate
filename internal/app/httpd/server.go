package httpd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"echo-boilerplate/internal/app/httpd/config"
	"echo-boilerplate/internal/app/httpd/handlers"
	"echo-boilerplate/internal/app/httpd/middleware"
	"echo-boilerplate/internal/pkg/logging"
)

func init() {
	c := config.NewConfig()
	c.BindFlags()
	lc := logging.Config{
		LogLevel:  viper.GetString("log-level"),
		LogOutput: viper.GetString("log-output"),
		LogWriter: viper.GetString("log-writer"),
	}
	err := logging.Init(lc)
	if err != nil {
		log.Panic().Msgf("Panic initializing logger: '%v'", err)
	}
}

// Start starts the echo HTTP server
func Start() {
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
