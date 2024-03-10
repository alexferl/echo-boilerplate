package cookie

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/util/jwt"
)

func TestNew(t *testing.T) {
	opts := &Options{
		Name:     "name",
		Value:    "value",
		Path:     "/path",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		MaxAge:   10,
	}
	c := New(opts)

	assert.Equal(t, opts.Name, c.Name)
	assert.Equal(t, opts.Value, c.Value)
	assert.Equal(t, opts.Path, c.Path)
	assert.Equal(t, opts.SameSite, c.SameSite)
	assert.Equal(t, opts.HttpOnly, c.HttpOnly)
	assert.Equal(t, opts.MaxAge, c.MaxAge)
}

func TestNewAccessToken(t *testing.T) {
	value := "access"
	cookie := NewAccessToken([]byte(value))

	assert.Equal(t, viper.GetString(config.JWTAccessTokenCookieName), cookie.Name)
	assert.Equal(t, value, cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	assert.Equal(t, int(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()), cookie.MaxAge)
}

func TestNewRefreshToken(t *testing.T) {
	value := "refresh"
	cookie := NewRefreshToken([]byte(value))

	assert.Equal(t, viper.GetString(config.JWTRefreshTokenCookieName), cookie.Name)
	assert.Equal(t, value, cookie.Value)
	assert.Equal(t, "/auth", cookie.Path)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	assert.Equal(t, int(viper.GetDuration(config.JWTRefreshTokenExpiry).Seconds()), cookie.MaxAge)
}

func TestNewCSRF(t *testing.T) {
	value := "csrf"
	cookie := NewCSRF([]byte(value))

	assert.Equal(t, viper.GetString(config.CSRFCookieName), cookie.Name)
	assert.Equal(t, value, cookie.Value)
	assert.Equal(t, "/", cookie.Path)
	assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	assert.Equal(t, int(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()), cookie.MaxAge)
}

func TestSetToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, resp)

	access, refresh, err := jwt.GenerateTokens("123", nil)
	assert.NoError(t, err)

	SetToken(ctx, access, refresh)

	accessCookie := resp.Result().Cookies()[0]
	refreshCookie := resp.Result().Cookies()[1]

	assert.Equal(t, string(access), accessCookie.Value)
	assert.Equal(t, int(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()), accessCookie.MaxAge)

	assert.Equal(t, string(refresh), refreshCookie.Value)
	assert.Equal(t, int(viper.GetDuration(config.JWTRefreshTokenExpiry).Seconds()), refreshCookie.MaxAge)
}

func TestSetExpiredToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, resp)

	SetExpiredToken(ctx)

	accessCookie := resp.Result().Cookies()[0]
	refreshCookie := resp.Result().Cookies()[1]

	assert.Equal(t, "", accessCookie.Value)
	assert.Equal(t, -1, accessCookie.MaxAge)

	assert.Equal(t, "", refreshCookie.Value)
	assert.Equal(t, -1, refreshCookie.MaxAge)
}
