package util

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/alexferl/echo-boilerplate/config"
)

type CookieOptions struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	SameSite http.SameSite
	HttpOnly bool
	MaxAge   int
}

func NewCookie(opts *CookieOptions) *http.Cookie {
	return &http.Cookie{
		Name:     opts.Name,
		Value:    opts.Value,
		Path:     opts.Path,
		Domain:   opts.Domain,
		SameSite: opts.SameSite,
		HttpOnly: opts.HttpOnly,
		Secure:   !(strings.ToUpper(viper.GetString(config.EnvName)) == "LOCAL"),
		MaxAge:   opts.MaxAge,
	}
}

func NewAccessTokenCookie(access []byte) *http.Cookie {
	opts := &CookieOptions{
		Name:     viper.GetString(config.JWTAccessTokenCookieName),
		Value:    string(access),
		Path:     "/",
		Domain:   viper.GetString(config.CookiesDomain),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()),
	}

	return NewCookie(opts)
}

func NewRefreshTokenCookie(refresh []byte) *http.Cookie {
	opts := &CookieOptions{
		Name:     viper.GetString(config.JWTRefreshTokenCookieName),
		Value:    string(refresh),
		Path:     "/auth",
		Domain:   viper.GetString(config.CookiesDomain),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		MaxAge:   int(viper.GetDuration(config.JWTRefreshTokenExpiry).Seconds()),
	}

	return NewCookie(opts)
}

func NewCSRFCookie(access []byte) *http.Cookie {
	opts := &CookieOptions{
		Name:     viper.GetString(config.CSRFCookieName),
		Value:    string(access),
		Path:     "/",
		Domain:   viper.GetString(config.CSRFCookieDomain),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()),
	}

	return NewCookie(opts)
}

func SetTokenCookies(c echo.Context, access []byte, refresh []byte) {
	c.SetCookie(NewAccessTokenCookie(access))
	c.SetCookie(NewRefreshTokenCookie(refresh))

	if viper.GetBool(config.CSRFEnabled) {
		if viper.GetString(config.CSRFSecretKey) == "" {
			log.Panic().Msg("CSRF secret key is unset!")
		}

		s := NewHMAC(access, []byte(viper.GetString(config.CSRFSecretKey)))

		c.SetCookie(NewCSRFCookie([]byte(s)))
	}
}

func SetExpiredTokenCookies(c echo.Context) {
	accessOpts := &CookieOptions{
		Name:     viper.GetString(config.JWTAccessTokenCookieName),
		Value:    "",
		Path:     "/",
		Domain:   viper.GetString(config.CookiesDomain),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}

	refreshOpts := &CookieOptions{
		Name:     viper.GetString(config.JWTRefreshTokenCookieName),
		Value:    "",
		Path:     "/auth",
		Domain:   viper.GetString(config.CookiesDomain),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		MaxAge:   -1,
	}

	c.SetCookie(NewCookie(accessOpts))
	c.SetCookie(NewCookie(refreshOpts))

	if viper.GetBool(config.CSRFEnabled) {
		csrfOpts := &CookieOptions{
			Name:     viper.GetString(config.CSRFCookieName),
			Value:    "",
			Path:     "/",
			Domain:   viper.GetString(config.CSRFCookieDomain),
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
		}
		c.SetCookie(NewCookie(csrfOpts))
	}
}
