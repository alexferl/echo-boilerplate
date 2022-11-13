package main

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	app "github.com/alexferl/echo-boilerplate"
	"github.com/alexferl/echo-boilerplate/config"
)

func main() {
	c := config.New()
	c.BindFlags()

	s := app.NewServer()

	log.Info().Msgf(
		"Starting %s on %s environment listening at %s",
		viper.GetString(config.AppName),
		strings.ToUpper(viper.GetString(config.EnvName)),
		fmt.Sprintf("%s:%d", viper.GetString(config.HTTPBindAddress), viper.GetInt(config.HTTPBindPort)),
	)

	s.Start()
}
