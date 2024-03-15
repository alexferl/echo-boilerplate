package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPersonalAccessToken(t *testing.T) {
	user := NewUser("test@email.com", "test")

	// expires_at not in the past
	_, err := NewPersonalAccessToken(user.Id, "My Token", time.Now().Format("2006-01-02"))
	assert.Error(t, err)
	assert.Equal(t, ErrExpiresAtPast, err)

	pat, err := NewPersonalAccessToken(user.Id, "My Token", time.Now().Add((7*24)*time.Hour).Format("2006-01-02"))
	assert.NoError(t, err)

	token := pat.Token
	create := pat.CreateResponse()
	assert.Equal(t, token, create.Token)

	resp := pat.Response()
	assert.Equal(t, pat.ExpiresAt, resp.ExpiresAt)

	err = pat.Encrypt()
	assert.NoError(t, err)

	err = pat.Validate(token)
	assert.NoError(t, err)
}

func TestPersonalAccessTokens(t *testing.T) {
	user := NewUser("test@email.com", "test")

	pat1, err := NewPersonalAccessToken(user.Id, "My Token1", time.Now().Add((7*24)*time.Hour).Format("2006-01-02"))
	assert.NoError(t, err)
	pat2, err := NewPersonalAccessToken(user.Id, "My Token1", time.Now().Add((7*24)*time.Hour).Format("2006-01-02"))
	assert.NoError(t, err)

	pats := PersonalAccessTokens{*pat1, *pat2}
	resp := pats.Response()
	assert.Len(t, resp.Tokens, 2)
}
