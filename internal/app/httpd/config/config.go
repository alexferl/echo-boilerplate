package config

import (
	"net"
	"net/http"

	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"

	"echo-boilerplate/internal/pkg/config"
)

// Config holds all configuration for our program
type Config struct {
	config.Config
	BindAddress         net.IP
	BindPort            uint
	CORS                middleware.CORSConfig
	CORSEnabled         bool
	GracefulTimeout     uint
	LogRequestsDisabled bool
}

// NewConfig creates a Config instance
func NewConfig() Config {
	cnf := Config{
		Config:      config.NewConfig(),
		BindAddress: net.ParseIP("127.0.0.1"),
		BindPort:    1323,
		CORS: middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
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
		CORSEnabled:         false,
		GracefulTimeout:     30,
		LogRequestsDisabled: false,
	}
	return cnf
}

// addFlags adds all the flags from the command line
func (cnf *Config) addFlags(fs *pflag.FlagSet) {
	fs.IPVar(&cnf.BindAddress, "bind-address", cnf.BindAddress, "The IP address to listen at.")
	fs.UintVar(&cnf.BindPort, "bind-port", cnf.BindPort, "The port to listen at.")
	fs.StringSliceVar(&cnf.CORS.AllowOrigins, "cors-allow-origins", cnf.CORS.AllowOrigins,
		"Indicates whether the response can be shared with requesting code from the given origin.")
	fs.StringSliceVar(&cnf.CORS.AllowMethods, "cors-allow-methods", cnf.CORS.AllowMethods,
		"Indicates which HTTP methods are allowed for cross-origin requests.")
	fs.StringSliceVar(&cnf.CORS.AllowHeaders, "cors-allow-headers", cnf.CORS.AllowHeaders,
		"Indicate which HTTP headers can be used during an actual request.")
	fs.BoolVar(&cnf.CORS.AllowCredentials, "cors-allow-credentials", cnf.CORS.AllowCredentials,
		"Tells browsers whether to expose the response to frontend JavaScript code when the request's credentials "+
			"mode (Request.credentials) is 'include'.")
	fs.StringSliceVar(&cnf.CORS.ExposeHeaders, "cors-expose-headers", cnf.CORS.ExposeHeaders,
		"Indicates which headers can be exposed as part of the response by listing their name.")
	fs.IntVar(&cnf.CORS.MaxAge, "cors-max-age", cnf.CORS.MaxAge,
		"Indicates how long the results of a preflight request can be cached.")
	fs.BoolVar(&cnf.CORSEnabled, "cors-enabled", cnf.CORSEnabled, "Enable cross-origin resource sharing.")
	fs.UintVar(&cnf.GracefulTimeout, "graceful-timeout", cnf.GracefulTimeout,
		"Timeout for graceful shutdown.")
	fs.BoolVar(&cnf.LogRequestsDisabled, "log-requests-disabled", cnf.LogRequestsDisabled,
		"Disables HTTP requests logging.")
}

// BindFlags normalizes and parses the command line flags
func (cnf *Config) BindFlags() {
	cnf.addFlags(pflag.CommandLine)
	cnf.Config.BindFlags()
}
