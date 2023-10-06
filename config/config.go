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

	Admin   *Admin
	OAuth2  *OAuth2
	JWT     *JWT
	Cookies *Cookies
	CSRF    *CSRF
	Casbin  *Casbin
	OpenAPI *OpenAPI
}

type Admin struct {
	Create   bool
	Email    string
	Username string
	Password string
}

type OAuth2 struct {
	ClientId     string
	ClientSecret string
}

type JWT struct {
	AccessTokenExpiry      time.Duration
	AccessTokenCookieName  string
	RefreshTokenExpiry     time.Duration
	RefreshTokenCookieName string
	PrivateKey             string
	Issuer                 string
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

type Casbin struct {
	Model  string
	Policy string
}

type OpenAPI struct {
	Schema string
}

// New creates a Config instance
func New() *Config {
	return &Config{
		Config:  libConfig.New("APP"),
		HTTP:    libHttp.DefaultConfig,
		Logging: libLog.DefaultConfig,
		MongoDB: libMongo.DefaultConfig,
		BaseURL: "http://localhost:1323",
		Admin: &Admin{
			Create:   false,
			Email:    "admin@example.com",
			Username: "admin",
			Password: "",
		},
		OAuth2: &OAuth2{
			ClientId:     "",
			ClientSecret: "",
		},
		JWT: &JWT{
			AccessTokenExpiry:      10 * time.Minute,
			AccessTokenCookieName:  "access_token",
			RefreshTokenExpiry:     (30 * 24) * time.Hour,
			RefreshTokenCookieName: "refresh_token",
			PrivateKey:             "./private-key.pem",
			Issuer:                 "http://localhost:1323",
		},
		Cookies: &Cookies{
			Enabled: false,
			Domain:  "",
		},
		CSRF: &CSRF{
			Enabled:      false,
			SecretKey:    "",
			CookieName:   "csrf_token",
			CookieDomain: "",
			HeaderName:   "X-CSRF-Token",
		},
		Casbin: &Casbin{
			Model:  "./casbin/model.conf",
			Policy: "./casbin/policy.csv",
		},
		OpenAPI: &OpenAPI{
			Schema: "./openapi/openapi.yaml",
		},
	}
}

const (
	AppName = libConfig.AppName
	EnvName = libConfig.EnvName

	HTTPBindAddress = libHttp.HTTPBindAddress
	HTTPBindPort    = libHttp.HTTPBindPort

	BaseURL = "base-url"

	AdminCreate   = "admin-create"
	AdminEmail    = "admin-email"
	AdminUsername = "admin-username"
	AdminPassword = "admin-password"

	OAuth2ClientId     = "oauth2-client-id"
	OAuth2ClientSecret = "oauth2-client-secret"

	JWTAccessTokenExpiry      = "jwt-access-token-expiry"
	JWTAccessTokenCookieName  = "jwt-access-token-cookie-name"
	JWTRefreshTokenExpiry     = "jwt-refresh-token-expiry"
	JWTRefreshTokenCookieName = "jwt-refresh-token-cookie-name"
	JWTPrivateKey             = "jwt-private-key"
	JWTIssuer                 = "jwt-issuer"

	CookiesEnabled = "cookies-enabled"
	CookiesDomain  = "cookies-domain"

	CSRFEnabled      = "csrf-enabled"
	CSRFSecretKey    = "csrf-secret-key"
	CSRFCookieName   = "csrf-cookie-name"
	CSRFCookieDomain = "csrf-cookie-domain"
	CSRFHeaderName   = "csrf-header-name"

	CasbinModel  = "casbin-model"
	CasbinPolicy = "casbin-policy"

	OpenAPISchema = "openapi-schema"
)

// addFlags adds all the flags from the command line
func (c *Config) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.BaseURL, BaseURL, c.BaseURL, "Base URL where the app will be served")

	fs.BoolVar(&c.Admin.Create, AdminCreate, c.Admin.Create, "Create admin")
	fs.StringVar(&c.Admin.Email, AdminEmail, c.Admin.Email, "Admin email")
	fs.StringVar(&c.Admin.Username, AdminUsername, c.Admin.Username, "Admin username")
	fs.StringVar(&c.Admin.Password, AdminPassword, c.Admin.Password, "Admin password")

	fs.StringVar(&c.OAuth2.ClientId, OAuth2ClientId, c.OAuth2.ClientId, "OAuth2 client id")
	fs.StringVar(&c.OAuth2.ClientSecret, OAuth2ClientSecret, c.OAuth2.ClientSecret, "OAuth2 client secret")

	fs.DurationVar(&c.JWT.AccessTokenExpiry, JWTAccessTokenExpiry, c.JWT.AccessTokenExpiry,
		"JWT access token expiry")
	fs.StringVar(&c.JWT.AccessTokenCookieName, JWTAccessTokenCookieName, c.JWT.AccessTokenCookieName,
		"JWT access token cookie name")
	fs.DurationVar(&c.JWT.RefreshTokenExpiry, JWTRefreshTokenExpiry, c.JWT.RefreshTokenExpiry,
		"JWT refresh token expiry")
	fs.StringVar(&c.JWT.RefreshTokenCookieName, JWTRefreshTokenCookieName, c.JWT.RefreshTokenCookieName,
		"JWT refresh token cookie name")
	fs.StringVar(&c.JWT.PrivateKey, JWTPrivateKey, c.JWT.PrivateKey, "JWT private key file path")
	fs.StringVar(&c.JWT.Issuer, JWTIssuer, c.JWT.Issuer, "JWT issuer")

	fs.BoolVar(&c.Cookies.Enabled, CookiesEnabled, c.Cookies.Enabled, "Send cookies with authentication requests")
	fs.StringVar(&c.Cookies.Domain, CookiesDomain, c.Cookies.Domain, "Cookies domain")

	fs.BoolVar(&c.CSRF.Enabled, CSRFEnabled, c.CSRF.Enabled, "CSRF enabled")
	fs.StringVar(&c.CSRF.SecretKey, CSRFSecretKey, c.CSRF.SecretKey, "CSRF secret used to hash the token")
	fs.StringVar(&c.CSRF.CookieName, CSRFCookieName, c.CSRF.CookieName, "CSRF cookie name")
	fs.StringVar(&c.CSRF.CookieDomain, CSRFCookieDomain, c.CSRF.CookieDomain, "CSRF cookie domain")
	fs.StringVar(&c.CSRF.HeaderName, CSRFHeaderName, c.CSRF.HeaderName, "CSRF header name")

	fs.StringVar(&c.Casbin.Model, CasbinModel, c.Casbin.Model, "Casbin model file")
	fs.StringVar(&c.Casbin.Policy, CasbinPolicy, c.Casbin.Policy, "Casbin policy file")

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

	err := c.Config.BindFlags()
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

	if viper.GetBool(AdminCreate) && viper.GetString(AdminPassword) == "" {
		log.Panic().Msg("Admin create: password is unset!")
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
