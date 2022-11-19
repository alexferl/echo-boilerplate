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

	accessToken, err := ParseToken(access)
	assert.NoError(t, err)
	assert.Equal(t, sub, accessToken.Subject())
	accessClaim, ok := accessToken.Get("claim")
	assert.True(t, ok)
	assert.Equal(t, claim, accessClaim)

	refreshToken, err := ParseToken(refresh)
	assert.NoError(t, err)
	assert.Equal(t, sub, refreshToken.Subject())
	_, ok = refreshToken.Get("claim")
	assert.False(t, ok)
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
		{"has role", map[string]any{"roles": []string{"user"}}, "user", true},
		{"invalid role", map[string]any{"roles": []string{"user"}}, "invalid", false},
		{"invalid roles key", map[string]any{"invalid": []string{"user"}}, "invalid", false},
		{"roles key not slice", map[string]any{"roles": "user"}, "user", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			access, _, err := GenerateTokens("123", tc.claims)
			assert.NoError(t, err)

			token, err := ParseToken(access)
			assert.NoError(t, err)

			assert.Equal(t, tc.hasRole, HasRole(token, tc.role))
		})
	}
}
