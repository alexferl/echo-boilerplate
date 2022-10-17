package app

import (
	xconfig "github.com/alexferl/golib/config"
	xhttp "github.com/alexferl/golib/http/config"
	xlog "github.com/alexferl/golib/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program
type Config struct {
	Config  *xconfig.Config
	Http    *xhttp.Config
	Logging *xlog.Config
}

// NewConfig creates a Config instance
func NewConfig() *Config {
	return &Config{
		Config:  xconfig.New("app"),
		Http:    xhttp.DefaultConfig,
		Logging: xlog.DefaultConfig,
	}
}

// addFlags adds all the flags from the command line
func (c *Config) addFlags(fs *pflag.FlagSet) {
	// add own flags
}

func (c *Config) BindFlags() {
	c.addFlags(pflag.CommandLine)
	c.Logging.BindFlags(pflag.CommandLine)
	c.Http.BindFlags(pflag.CommandLine)

	err := c.Config.BindFlags()
	if err != nil {
		panic(err)
	}

	err = xlog.New(&xlog.Config{
		LogLevel:  viper.GetString("log-level"),
		LogOutput: viper.GetString("log-output"),
		LogWriter: viper.GetString("log-writer"),
	})
	if err != nil {
		panic(err)
	}
}
