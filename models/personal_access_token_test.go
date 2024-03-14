package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPersonalAccessToken(t *testing.T) {
	user := NewUser("test@email.com", "test")

	// expires_at not in future
	_, err := NewPersonalAccessToken(user.Id, "My Token", time.Now().Format("2006-01-02"))
	assert.Error(t, err)

	pat, err := NewPersonalAccessToken(user.Id, "My Token", time.Now().Add((7*24)*time.Hour).Format("2006-01-02"))
	assert.NoError(t, err)

	err = pat.Encrypt()
	assert.NoError(t, err)
}
