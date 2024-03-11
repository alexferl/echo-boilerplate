package config

import (
	"fmt"
	"time"

	libConfig "github.com/alexferl/golib/config"
	libMongo "github.com/alexferl/golib/database/mongodb"
	libHttp "github.com/alexferl/golib/http/api/config"
	libLog "github.com/alexferl/golib/log"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program
type Config struct {
	Config  *libConfig.Config
	HTTP    *libHttp.Config
	Logging *libLog.Config
	MongoDB *libMongo.Config

	BaseURL string

	Casbin       *Casbin
	Cookies      *Cookies
	CSRF         *CSRF
	JWT          *JWT
	OAuth2       *OAuth2
	OAuth2Google *OAuth2Google
	OpenAPI      *OpenAPI
}

type Casbin struct {
	Model  string
	Policy string
}

type Cookies struct {
	Enabled bool
	Domain  string
}

type CSRF struct {
	Enabled      bool
	SecretKey    string
	CookieName   string
	CookieDomain string
	HeaderName   string
}

type JWT struct {
	AccessTokenExpiry      time.Duration
	AccessTokenCookieName  string
	RefreshTokenExpiry     time.Duration
	RefreshTokenCookieName string
	PrivateKey             string
	Issuer                 string
}

type OAuth2 struct {
	Providers []string
}

type OAuth2Google struct {
	ClientId     string
	ClientSecret string
}

type OpenAPI struct {
	Schema string
}

// New creates a Config instance
func New() *Config {
	c := &Config{
		Config:  libConfig.New("APP"),
		HTTP:    libHttp.DefaultConfig,
		Logging: libLog.DefaultConfig,
		MongoDB: libMongo.DefaultConfig,
		BaseURL: "http://localhost:1323",
		Casbin: &Casbin{
			Model:  "./casbin/model.conf",
			Policy: "./casbin/policy.csv",
		},
		Cookies: &Cookies{
			Enabled: false,
			Domain:  "",
		},
		CSRF: &CSRF{
			Enabled:      false,
			CookieDomain: "",
			CookieName:   "csrf_token",
			HeaderName:   "X-CSRF-Token",
			SecretKey:    "",
		},
		JWT: &JWT{
			AccessTokenCookieName:  "access_token",
			AccessTokenExpiry:      60 * time.Minute,
			PrivateKey:             "./private-key.pem",
			RefreshTokenCookieName: "refresh_token",
			RefreshTokenExpiry:     (30 * 24) * time.Hour,
		},
		OAuth2: &OAuth2{
			Providers: []string{""},
		},
		OAuth2Google: &OAuth2Google{
			ClientId:     "",
			ClientSecret: "",
		},
		OpenAPI: &OpenAPI{
			Schema: "./openapi/openapi.yaml",
		},
	}
	c.JWT.Issuer = c.BaseURL
	return c
}

const (
	AppName = libConfig.AppName
	EnvName = libConfig.EnvName

	HTTPBindAddress = libHttp.HTTPBindAddress
	HTTPBindPort    = libHttp.HTTPBindPort

	BaseURL = "base-url"

	CasbinModel  = "casbin-model"
	CasbinPolicy = "casbin-policy"

	CookiesEnabled = "cookies-enabled"
	CookiesDomain  = "cookies-domain"

	CSRFEnabled      = "csrf-enabled"
	CSRFCookieDomain = "csrf-cookie-domain"
	CSRFCookieName   = "csrf-cookie-name"
	CSRFHeaderName   = "csrf-header-name"
	CSRFSecretKey    = "csrf-secret-key"

	JWTAccessTokenCookieName  = "jwt-access-token-cookie-name"
	JWTAccessTokenExpiry      = "jwt-access-token-expiry"
	JWTIssuer                 = "jwt-issuer"
	JWTPrivateKey             = "jwt-private-key"
	JWTRefreshTokenCookieName = "jwt-refresh-token-cookie-name"
	JWTRefreshTokenExpiry     = "jwt-refresh-token-expiry"

	OAuth2Providers = "oauth2-providers"

	OAuth2GoogleClientId     = "oauth2-google-client-id"
	OAuth2GoogleClientSecret = "oauth2-google-client-secret"

	OpenAPISchema = "openapi-schema"
)

// addFlags adds all the flags from the command line
func (c *Config) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.BaseURL, BaseURL, c.BaseURL, "Base URL where the app will be served")

	fs.StringVar(&c.Casbin.Model, CasbinModel, c.Casbin.Model, "Casbin model file")
	fs.StringVar(&c.Casbin.Policy, CasbinPolicy, c.Casbin.Policy, "Casbin policy file")

	fs.BoolVar(&c.Cookies.Enabled, CookiesEnabled, c.Cookies.Enabled, "Send cookies with authentication requests")
	fs.StringVar(&c.Cookies.Domain, CookiesDomain, c.Cookies.Domain, "Cookies domain")

	fs.BoolVar(&c.CSRF.Enabled, CSRFEnabled, c.CSRF.Enabled, "CSRF enabled")
	fs.StringVar(&c.CSRF.SecretKey, CSRFSecretKey, c.CSRF.SecretKey, "CSRF secret used to hash the token")
	fs.StringVar(&c.CSRF.CookieName, CSRFCookieName, c.CSRF.CookieName, "CSRF cookie name")
	fs.StringVar(&c.CSRF.CookieDomain, CSRFCookieDomain, c.CSRF.CookieDomain, "CSRF cookie domain")
	fs.StringVar(&c.CSRF.HeaderName, CSRFHeaderName, c.CSRF.HeaderName, "CSRF header name")

	fs.StringVar(&c.JWT.AccessTokenCookieName, JWTAccessTokenCookieName, c.JWT.AccessTokenCookieName,
		"JWT access token cookie name")
	fs.DurationVar(&c.JWT.AccessTokenExpiry, JWTAccessTokenExpiry, c.JWT.AccessTokenExpiry,
		"JWT access token expiry")
	fs.StringVar(&c.JWT.Issuer, JWTIssuer, c.JWT.Issuer, "JWT issuer")
	fs.StringVar(&c.JWT.PrivateKey, JWTPrivateKey, c.JWT.PrivateKey, "JWT private key file path")
	fs.StringVar(&c.JWT.RefreshTokenCookieName, JWTRefreshTokenCookieName, c.JWT.RefreshTokenCookieName,
		"JWT refresh token cookie name")
	fs.DurationVar(&c.JWT.RefreshTokenExpiry, JWTRefreshTokenExpiry, c.JWT.RefreshTokenExpiry,
		"JWT refresh token expiry")

	fs.StringSliceVar(&c.OAuth2.Providers, OAuth2Providers, c.OAuth2.Providers, "OAuth2 providers")

	fs.StringVar(&c.OAuth2Google.ClientId, OAuth2GoogleClientId, c.OAuth2Google.ClientId, "OAuth2 Google client id")
	fs.StringVar(&c.OAuth2Google.ClientSecret, OAuth2GoogleClientSecret, c.OAuth2Google.ClientSecret, "OAuth2 Google client secret")

	fs.StringVar(&c.OpenAPI.Schema, OpenAPISchema, c.OpenAPI.Schema, "OpenAPI schema file")
}

func (c *Config) BindFlags() {
	if pflag.Parsed() {
		return
	}

	c.addFlags(pflag.CommandLine)
	c.Logging.BindFlags(pflag.CommandLine)
	c.HTTP.BindFlags(pflag.CommandLine)
	c.MongoDB.BindFlags(pflag.CommandLine)

	err := c.Config.BindFlagsWithConfigPaths()
	if err != nil {
		panic(fmt.Errorf("failed binding flags: %v", err))
	}

	err = libLog.New(&libLog.Config{
		LogLevel:  viper.GetString(libLog.LogLevel),
		LogOutput: viper.GetString(libLog.LogOutput),
		LogWriter: viper.GetString(libLog.LogWriter),
	})
	if err != nil {
		panic(fmt.Errorf("failed creating logger: %v", err))
	}

	if viper.GetBool(CSRFEnabled) && viper.GetString(CSRFSecretKey) == "" {
		log.Panic().Msg("CSRF: secret key is unset!")
	}

	if viper.GetBool(libHttp.HTTPCORSEnabled) {
		for _, origin := range viper.GetStringSlice(libHttp.HTTPCORSAllowOrigins) {
			if origin == "*" {
				log.Warn().Msg("CORS: using '*' in Access-Control-Allow-Origin is potentially unsafe!")
			}

			if origin == "null" {
				log.Warn().Msg("CORS: using 'null' in Access-Control-Allow-Origin is unsafe and should not be used!")
			}

		}
	}
}
