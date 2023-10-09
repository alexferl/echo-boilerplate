package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	pwd := "s3cret"
	enc, err := HashPassword([]byte(pwd))
	assert.NoError(t, err)

	err = VerifyPassword([]byte(enc), []byte(pwd))
	assert.NoError(t, err)
}

func TestWrongPassword(t *testing.T) {
	pwd := "s3cret"
	enc, err := HashPassword([]byte(pwd))
	assert.NoError(t, err)

	err = VerifyPassword([]byte(enc), []byte("wrong"))
	assert.Error(t, err)
}
