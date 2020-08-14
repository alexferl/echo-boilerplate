package logging

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes the logger based on the config
func Init(config *Config) {
	logLevel := strings.ToLower(config.LogLevel)
	logOutput := strings.ToLower(config.LogOutput)
	logWriter := strings.ToLower(config.LogWriter)

	var f *os.File
	switch logOutput {
	case "stdout":
		f = os.Stdout
	case "stderr":
		f = os.Stderr
	default:
		log.Warn().Msgf("Unknown log-output '%s', falling back to 'stdout'", logWriter)
		f = os.Stdout
	}

	logger := zerolog.New(f)

	switch strings.ToLower(logWriter) {
	case "console":
		logger = log.Output(zerolog.ConsoleWriter{Out: f})
	case "json":
		break
	default:
		log.Warn().Msgf("Unknown log-writer '%s', falling back to 'console'", logWriter)
		logger = log.Output(zerolog.ConsoleWriter{Out: f})
	}

	log.Logger = logger.With().Timestamp().Caller().Logger()

	switch strings.ToLower(logLevel) {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		log.Warn().Msgf("Unknown log-level '%s', falling back to '%s'", logLevel, zerolog.InfoLevel)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
