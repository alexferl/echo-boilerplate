package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/alexferl/echo-boilerplate/config"
)

func TestNewCookie(t *testing.T) {
	opts := &CookieOptions{
		Name:     "name",
		Value:    "value",
		Path:     "/path",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		MaxAge:   10,
	}
	c := NewCookie(opts)

	assert.Equal(t, opts.Name, c.Name)
	assert.Equal(t, opts.Value, c.Value)
	assert.Equal(t, opts.Path, c.Path)
	assert.Equal(t, opts.SameSite, c.SameSite)
	assert.Equal(t, opts.HttpOnly, c.HttpOnly)
	assert.Equal(t, opts.MaxAge, c.MaxAge)
}

func TestNewAccessTokenCookie(t *testing.T) {
	value := "access"
	c := NewAccessTokenCookie([]byte(value))

	assert.Equal(t, "access_token", c.Name)
	assert.Equal(t, value, c.Value)
	assert.Equal(t, "/", c.Path)
	assert.Equal(t, http.SameSiteStrictMode, c.SameSite)
	assert.Equal(t, viper.GetInt(config.JWTAccessTokenExpiry), c.MaxAge)
}

func TestNewRefreshTokenCookie(t *testing.T) {
	value := "refresh"
	c := NewRefreshTokenCookie([]byte(value))

	assert.Equal(t, "refresh_token", c.Name)
	assert.Equal(t, value, c.Value)
	assert.Equal(t, "/auth", c.Path)
	assert.Equal(t, http.SameSiteStrictMode, c.SameSite)
	assert.Equal(t, viper.GetInt(config.JWTRefreshTokenExpiry), c.MaxAge)
}

func TestSetTokenCookies(t *testing.T) {
	c := config.New()
	c.BindFlags()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, resp)

	access, refresh, err := GenerateTokens("123", nil)
	assert.NoError(t, err)

	SetTokenCookies(ctx, access, refresh)

	accessCookie := resp.Result().Cookies()[0]
	refreshCookie := resp.Result().Cookies()[1]

	assert.Equal(t, string(access), accessCookie.Value)
	assert.Equal(t, int(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()), accessCookie.MaxAge)

	assert.Equal(t, string(refresh), refreshCookie.Value)
	assert.Equal(t, int(viper.GetDuration(config.JWTRefreshTokenExpiry).Seconds()), refreshCookie.MaxAge)
}

func TestSetExpiredTokenCookies(t *testing.T) {
	c := config.New()
	c.BindFlags()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, resp)

	SetExpiredTokenCookies(ctx)

	accessCookie := resp.Result().Cookies()[0]
	refreshCookie := resp.Result().Cookies()[1]

	assert.Equal(t, "", accessCookie.Value)
	assert.Equal(t, -1, accessCookie.MaxAge)

	assert.Equal(t, "", refreshCookie.Value)
	assert.Equal(t, -1, refreshCookie.MaxAge)
}
