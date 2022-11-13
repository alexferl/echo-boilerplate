package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	s, err := GenerateRandomString(10)

	assert.NoError(t, err)
	assert.Equal(t, 16, len(s))
}
