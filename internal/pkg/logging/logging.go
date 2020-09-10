package logging

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes the logger based on the config
func Init(config Config) error {
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
		return errors.New(fmt.Sprintf("Unknown log-output '%s'", logOutput))
	}

	logger := zerolog.New(f)

	switch strings.ToLower(logWriter) {
	case "console":
		logger = log.Output(zerolog.ConsoleWriter{Out: f})
	case "json":
		break
	default:
		return errors.New(fmt.Sprintf("Unknown log-writer '%s'", logWriter))
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
		return errors.New(fmt.Sprintf("Unknown log-level '%s'", logLevel))
	}

	return nil
}
