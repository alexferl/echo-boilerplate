package logging

import "github.com/spf13/pflag"

type Config struct {
	LogLevel  string
	LogOutput string
	LogWriter string
}

// NewConfig creates a Config instance
func NewConfig() *Config {
	cnf := Config{
		LogOutput: "stdout",
		LogWriter: "json",
		LogLevel:  "info",
	}
	return &cnf
}

// AddFlags adds all the flags from the command line
func (c *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.LogOutput, "log-output", c.LogOutput, "The output to write to. "+
		"'stdout' means log to stdout, 'stderr' means log to stderr.")
	fs.StringVar(&c.LogWriter, "log-writer", c.LogWriter,
		"The log writer. Valid writers are: 'console' and 'json'.")
	fs.StringVar(&c.LogLevel, "log-level", c.LogLevel, "The granularity of log outputs. "+
		"Valid log levels: 'panic', 'fatal', 'error', 'warn', 'info', 'debug' and 'trace'.")
}
