package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alexferl/echo-boilerplate/util/jwt"
)

func TestPersonalAccessToken(t *testing.T) {
	user := NewUser("test@email.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	token, err := jwt.ParseEncoded(access)
	assert.NoError(t, err)

	// expires_at not in future
	_, err = NewPersonalAccessToken(token, "My Token", time.Now().Format("2006-01-02"))
	assert.Error(t, err)

	pat, err := NewPersonalAccessToken(token, "My Token", time.Now().Add((7*24)*time.Hour).Format("2006-01-02"))
	assert.NoError(t, err)

	err = pat.Encrypt()
	assert.NoError(t, err)
}
