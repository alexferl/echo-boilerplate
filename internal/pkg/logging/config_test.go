package logging

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	c := NewConfig()

	assert.Equal(t, DefaultConfig.LogLevel, c.LogLevel)
	assert.Equal(t, DefaultConfig.LogOutput, c.LogOutput)
	assert.Equal(t, DefaultConfig.LogWriter, c.LogWriter)
}

func TestNewConfigWithConfig(t *testing.T) {
	c1 := Config{LogLevel: "warn"}
	cc1 := NewConfigWithConfig(c1)
	assert.Equal(t, c1.LogLevel, cc1.LogLevel)

	c2 := Config{LogOutput: "stderr"}
	cc2 := NewConfigWithConfig(c2)
	assert.Equal(t, c2.LogOutput, cc2.LogOutput)

	c3 := Config{LogWriter: "console"}
	cc3 := NewConfigWithConfig(c3)
	assert.Equal(t, c3.LogWriter, cc3.LogWriter)
}

func TestConfig_AddFlags(t *testing.T) {
	c := NewConfig()
	fs := &pflag.FlagSet{}
	c.AddFlags(fs)

	level, _ := fs.GetString("log-level")
	out, _ := fs.GetString("log-output")
	writer, _ := fs.GetString("log-writer")
	assert.Equal(t, DefaultConfig.LogLevel, level)
	assert.Equal(t, DefaultConfig.LogOutput, out)
	assert.Equal(t, DefaultConfig.LogWriter, writer)
}
