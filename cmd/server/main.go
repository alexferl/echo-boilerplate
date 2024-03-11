package main

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/server"
)

func main() {
	c := config.New()
	c.BindFlags()

	s := server.New()

	log.Info().Msgf(
		"Starting %s on %s environment listening at http://%s",
		viper.GetString(config.AppName),
		strings.ToUpper(viper.GetString(config.EnvName)),
		fmt.Sprintf("%s:%d", viper.GetString(config.HTTPBindAddress), viper.GetInt(config.HTTPBindPort)),
	)

	s.Start()
}
