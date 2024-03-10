package cookie

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/util/hash"
)

type Options struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	SameSite http.SameSite
	HttpOnly bool
	MaxAge   int
}

func New(opts *Options) *http.Cookie {
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

func NewAccessToken(access []byte) *http.Cookie {
	opts := &Options{
		Name:     viper.GetString(config.JWTAccessTokenCookieName),
		Value:    string(access),
		Path:     "/",
		Domain:   viper.GetString(config.CookiesDomain),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()),
	}

	return New(opts)
}

func NewRefreshToken(refresh []byte) *http.Cookie {
	opts := &Options{
		Name:     viper.GetString(config.JWTRefreshTokenCookieName),
		Value:    string(refresh),
		Path:     "/auth",
		Domain:   viper.GetString(config.CookiesDomain),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		MaxAge:   int(viper.GetDuration(config.JWTRefreshTokenExpiry).Seconds()),
	}

	return New(opts)
}

func NewCSRF(access []byte) *http.Cookie {
	opts := &Options{
		Name:     viper.GetString(config.CSRFCookieName),
		Value:    string(access),
		Path:     "/",
		Domain:   viper.GetString(config.CSRFCookieDomain),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()),
	}

	return New(opts)
}

func SetToken(c echo.Context, access []byte, refresh []byte) {
	c.SetCookie(NewAccessToken(access))
	c.SetCookie(NewRefreshToken(refresh))

	if viper.GetBool(config.CSRFEnabled) {
		s := hash.NewHMAC(access, []byte(viper.GetString(config.CSRFSecretKey)))
		c.SetCookie(NewCSRF([]byte(s)))
	}
}

func SetExpiredToken(c echo.Context) {
	accessOpts := &Options{
		Name:     viper.GetString(config.JWTAccessTokenCookieName),
		Value:    "",
		Path:     "/",
		Domain:   viper.GetString(config.CookiesDomain),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}

	refreshOpts := &Options{
		Name:     viper.GetString(config.JWTRefreshTokenCookieName),
		Value:    "",
		Path:     "/auth",
		Domain:   viper.GetString(config.CookiesDomain),
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		MaxAge:   -1,
	}

	c.SetCookie(New(accessOpts))
	c.SetCookie(New(refreshOpts))

	if viper.GetBool(config.CSRFEnabled) {
		csrfOpts := &Options{
			Name:     viper.GetString(config.CSRFCookieName),
			Value:    "",
			Path:     "/",
			Domain:   viper.GetString(config.CSRFCookieDomain),
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
		}
		c.SetCookie(New(csrfOpts))
	}
}
