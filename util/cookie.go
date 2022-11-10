package util

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"github.com/alexferl/echo-boilerplate/config"
)

type CookieOptions struct {
	Name     string
	Value    string
	Path     string
	SameSite http.SameSite
	HttpOnly bool
	MaxAge   int
}

func NewCookie(opts *CookieOptions) *http.Cookie {
	var secure bool
	if !(strings.ToUpper(viper.GetString(config.EnvName)) == "LOCAL") {
		secure = true
	}

	return &http.Cookie{
		Name:     opts.Name,
		Value:    opts.Value,
		Path:     opts.Path,
		SameSite: opts.SameSite,
		HttpOnly: opts.HttpOnly,
		Secure:   secure,
		MaxAge:   opts.MaxAge,
	}
}

func NewAccessTokenCookie(access string) *http.Cookie {
	opts := &CookieOptions{
		Name:     "access_token",
		Value:    access,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   viper.GetInt(config.JWTAccessTokenExpiry),
	}

	return NewCookie(opts)
}

func NewRefreshTokenCookie(refresh string) *http.Cookie {
	opts := &CookieOptions{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/auth",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		MaxAge:   viper.GetInt(config.JWTRefreshTokenExpiry),
	}

	return NewCookie(opts)
}

func SetTokenCookies(c echo.Context, access string, refresh string) {
	c.SetCookie(NewAccessTokenCookie(access))
	c.SetCookie(NewRefreshTokenCookie(refresh))
}

func SetExpiredTokenCookies(c echo.Context) {
	accessOpts := &CookieOptions{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}

	refreshOpts := &CookieOptions{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		MaxAge:   -1,
	}

	c.SetCookie(NewCookie(accessOpts))
	c.SetCookie(NewCookie(refreshOpts))
}
