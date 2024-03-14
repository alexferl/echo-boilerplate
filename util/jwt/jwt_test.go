package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexferl/echo-boilerplate/config"
)

func TestGenerateTokens(t *testing.T) {
	c := config.New()
	c.BindFlags()

	_, _, err := GenerateTokens("123", nil)
	assert.NoError(t, err)
}

func TestParseEncoded(t *testing.T) {
	c := config.New()
	c.BindFlags()

	sub := "123"
	claim := "mine"

	access, refresh, err := GenerateTokens(sub, map[string]any{"claim": claim})
	assert.NoError(t, err)

	accessToken, err := ParseEncoded(access)
	assert.NoError(t, err)
	assert.Equal(t, sub, accessToken.Subject())
	accessClaim, ok := accessToken.Get("claim")
	assert.True(t, ok)
	assert.Equal(t, claim, accessClaim)

	refreshToken, err := ParseEncoded(refresh)
	assert.NoError(t, err)
	assert.Equal(t, sub, refreshToken.Subject())
	_, ok = refreshToken.Get("claim")
	assert.False(t, ok)
}
