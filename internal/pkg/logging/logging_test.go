package logging

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	var tests = []struct {
		config Config
		fail   bool
	}{
		{DefaultConfig, false},
		{Config{LogLevel: "warn", LogOutput: "stdout", LogWriter: "json"}, false},
		{Config{LogLevel: "info", LogOutput: "stderr", LogWriter: "json"}, false},
		{Config{LogLevel: "info", LogOutput: "stdout", LogWriter: "console"}, false},
		{Config{LogLevel: "panic", LogOutput: "stdout", LogWriter: "json"}, false},
		{Config{LogLevel: "fatal", LogOutput: "stdout", LogWriter: "json"}, false},
		{Config{LogLevel: "error", LogOutput: "stdout", LogWriter: "json"}, false},
		{Config{LogLevel: "warn", LogOutput: "stdout", LogWriter: "json"}, false},
		{Config{LogLevel: "debug", LogOutput: "stdout", LogWriter: "json"}, false},
		{Config{LogLevel: "trace", LogOutput: "stdout", LogWriter: "json"}, false},
		{Config{LogLevel: "wrong"}, true},
		{Config{LogLevel: "info", LogOutput: "wrong"}, true},
		{Config{LogLevel: "info", LogOutput: "stdout", LogWriter: "wrong"}, true},
	}

	for _, tt := range tests {
		err := Init(tt.config)
		if !tt.fail {
			assert.NoError(t, err)
			assert.Equal(t, tt.config.LogLevel, zerolog.GlobalLevel().String())
		} else {
			assert.Error(t, err)
		}
	}
}
