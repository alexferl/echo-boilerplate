package util

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexferl/echo-boilerplate/config"
	_ "github.com/alexferl/echo-boilerplate/testing"
)

func TestGenerateTokens(t *testing.T) {
	c := config.New()
	c.BindFlags()

	_, _, err := GenerateTokens("123", nil)
	assert.NoError(t, err)
}

func TestParseToken(t *testing.T) {
	c := config.New()
	c.BindFlags()

	sub := "123"
	claim := "mine"

	access, refresh, err := GenerateTokens(sub, map[string]any{"claim": claim})
	assert.NoError(t, err)

	accessToken, err := parseToken(access)
	assert.NoError(t, err)
	assert.Equal(t, sub, accessToken.Subject())
	accessClaim, ok := accessToken.Get("claim")
	assert.True(t, ok)
	assert.Equal(t, claim, accessClaim)

	refreshToken, err := parseToken(refresh)
	assert.NoError(t, err)
	assert.Equal(t, sub, refreshToken.Subject())
	_, ok = refreshToken.Get("claim")
	assert.False(t, ok)
}

func TestHashToken(t *testing.T) {
	c := config.New()
	c.BindFlags()

	access, _, err := GenerateTokens("123", nil)
	assert.NoError(t, err)

	token, err := parseToken(access)
	assert.NoError(t, err)

	_, err = HashToken(token)
	assert.NoError(t, err)
}

func TestHasRole(t *testing.T) {
	c := config.New()
	c.BindFlags()

	testCases := []struct {
		name    string
		claims  map[string]any
		role    string
		hasRole bool
	}{
		{"", map[string]any{"roles": []string{"user"}}, "user", true},
		{"", map[string]any{"roles": []string{"user"}}, "invalid", false},
		{"", map[string]any{"invalid": []string{"user"}}, "invalid", false},
		{"", map[string]any{"roles": "user"}, "user", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			access, _, err := GenerateTokens("123", tc.claims)
			assert.NoError(t, err)

			token, err := parseToken(access)
			assert.NoError(t, err)

			assert.Equal(t, tc.hasRole, HasRole(token, tc.role))
		})
	}
}