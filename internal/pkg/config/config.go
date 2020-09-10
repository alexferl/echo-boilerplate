package config

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"echo-boilerplate/internal/pkg/logging"
)

// Config holds all global configuration for our program
type Config struct {
	AppName      string
	EnvName      string
	EnvVarPrefix string
	Logging      logging.Config
}

// NewConfig creates a Config instance
func NewConfig() Config {
	cnf := Config{
		Logging:      logging.NewConfig(),
		AppName:      "app",
		EnvName:      "local",
		EnvVarPrefix: "app",
	}
	return cnf
}

// addFlags adds all the flags from the command line
func (cnf *Config) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cnf.AppName, "app-name", cnf.AppName, "The name of the application.")
	fs.StringVar(&cnf.EnvName, "env-name", cnf.EnvName, "The environment of the application. "+
		"Used to load the right config file.")
	fs.StringVar(&cnf.EnvVarPrefix, "env-var-prefix", cnf.EnvVarPrefix,
		"Used to prefix environment variables.")
}

// wordSepNormalizeFunc changes all flags that contain "_" separators
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}

// BindFlags normalizes and parses the command line flags
func (cnf *Config) BindFlags() {
	cnf.Logging.AddFlags(pflag.CommandLine)

	cnf.addFlags(pflag.CommandLine)
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatal().Msgf("Error binding flags: '%v'", err)
	}

	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	pflag.Parse()

	n := viper.GetString("app-name")
	if len(n) < 1 {
		log.Fatal().Msgf("Application name cannot be empty!")
	}

	viper.SetEnvPrefix(n)
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	configName := fmt.Sprintf("config.%s", strings.ToLower(viper.GetString("env-name")))
	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/configs")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Error().Msgf("Config file not found: '%v'", err)
		} else {
			log.Fatal().Msgf("Couldn't load config file: '%v'", err)
		}
	}
}
