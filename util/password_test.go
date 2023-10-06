package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	pwd := "s3cret"
	enc, err := HashPassword(pwd)
	assert.NoError(t, err)

	err = VerifyPassword(enc, pwd)
	assert.NoError(t, err)
}
