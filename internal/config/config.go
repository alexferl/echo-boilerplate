package config

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program
type Config struct {
	AppName             string
	EnvName             string
	BindAddress         net.IP
	BindPort            uint
	GracefulTimeout     uint
	LogFile             string
	LogFormat           string
	LogLevel            string
	LogRequestsDisabled bool
	CORS                middleware.CORSConfig
}

// NewConfig creates a Config instance
func NewConfig() *Config {
	cnf := Config{
		AppName:             "app",
		EnvName:             "local",
		BindAddress:         net.ParseIP("127.0.0.1"),
		BindPort:            1323,
		GracefulTimeout:     30,
		LogFile:             "stdout",
		LogFormat:           "text",
		LogLevel:            "info",
		LogRequestsDisabled: false,
		CORS: middleware.CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPut,
				http.MethodPatch,
				http.MethodPost,
				http.MethodDelete,
			},
			AllowHeaders:     []string{},
			AllowCredentials: false,
			ExposeHeaders:    []string{},
			MaxAge:           0,
		},
	}
	return &cnf
}

// addFlags adds all the flags from the command line
func (cnf *Config) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cnf.AppName, "app-name", cnf.AppName, "The name of the application." +
		"Used to prefix environment variables.")
	fs.StringVar(&cnf.EnvName, "env-name", cnf.EnvName, "The environment of the application. " +
		"Used to load the right config file.")
	fs.IPVar(&cnf.BindAddress, "bind-address", cnf.BindAddress, "The IP address to listen at.")
	fs.UintVar(&cnf.BindPort, "bind-port", cnf.BindPort, "The port to listen at.")
	fs.UintVar(&cnf.GracefulTimeout, "graceful-timeout", cnf.GracefulTimeout,
		"Timeout for graceful shutdown.")
	fs.StringVar(&cnf.LogFile, "log-file", cnf.LogFile, "The log file to write to. "+
		"'stdout' means log to stdout, 'stderr' means log to stderr and 'null' means discard log messages.")
	fs.StringVar(&cnf.LogFormat, "log-format", cnf.LogFormat,
		"The log format. Valid format values are: text, json.")
	fs.StringVar(&cnf.LogLevel, "log-level", cnf.LogLevel, "The granularity of log outputs. "+
		"Valid log levels: debug, info, warning, error and critical.")
	fs.BoolVar(&cnf.LogRequestsDisabled, "log-requests-disabled", cnf.LogRequestsDisabled,
		"Disables HTTP requests logging.")
	fs.StringSliceVar(&cnf.CORS.AllowOrigins, "cors-allow-origins", cnf.CORS.AllowOrigins,
		"Indicates whether the response can be shared with requesting code from the given origin")
	fs.StringSliceVar(&cnf.CORS.AllowMethods, "cors-allow-methods", cnf.CORS.AllowMethods,
		"Indicates which HTTP methods are allowed for cross-origin requests.")
	fs.StringSliceVar(&cnf.CORS.AllowHeaders, "cors-allow-headers", cnf.CORS.AllowHeaders,
		"Indicate which HTTP headers can be used during an actual request.")
	fs.BoolVar(&cnf.CORS.AllowCredentials, "cors-allow-credentials", cnf.CORS.AllowCredentials,
		"Tells browsers whether to expose the response to frontend JavaScript code when the request's credentials " +
		"mode (Request.credentials) is 'include'.")
	fs.StringSliceVar(&cnf.CORS.ExposeHeaders, "cors-expose-headers", cnf.CORS.ExposeHeaders,
		"Indicates which headers can be exposed as part of the response by listing their name.")
	fs.IntVar(&cnf.CORS.MaxAge, "cors-max-age", cnf.CORS.MaxAge,
		"Indicates how long the results of a preflight request can be cached.")
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
	cnf.addFlags(pflag.CommandLine)
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		m := fmt.Sprintf("Error binding flags: %v", err)
		logrus.Panic(m)
		panic(m)
	}

	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	pflag.Parse()

	n := viper.GetString("app-name")
	if len(n) < 1 {
		m := fmt.Sprint("Application name cannot be empty!")
		logrus.Panic(m)
		panic(m)
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

	err = viper.ReadInConfig()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Errorf("Config file not found: %v", err)
		} else {
			logrus.Panic("Couldn't load config file: %v", err)
		}
	}
}
