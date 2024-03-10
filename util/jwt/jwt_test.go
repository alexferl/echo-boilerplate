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

func TestHasRoles(t *testing.T) {
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

			token, err := ParseEncoded(access)
			assert.NoError(t, err)

			assert.Equal(t, tc.hasRole, HasRoles(token, tc.role))
		})
	}
}

func TestGetRoles(t *testing.T) {
	c := config.New()
	c.BindFlags()

	testCases := []struct {
		name  string
		roles []string
	}{
		{"no role", []string{""}},
		{"user role", []string{"user"}},
		{"many roles", []string{"user", "admin", "super"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			access, _, err := GenerateTokens("123", map[string]any{"roles": tc.roles})
			assert.NoError(t, err)

			token, err := ParseEncoded(access)
			assert.NoError(t, err)

			assert.Equal(t, tc.roles, GetRoles(token))
		})
	}
}
