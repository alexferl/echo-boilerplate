package config

import (
	"time"

	libconfig "github.com/alexferl/golib/config"
	libhttp "github.com/alexferl/golib/http/config"
	liblog "github.com/alexferl/golib/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds all configuration for our program
type Config struct {
	Config  *libconfig.Config
	HTTP    *libhttp.Config
	Logging *liblog.Config
	BaseUrl string
	Admin   *Admin
	OAuth2  *OAuth2
	JWT     *JWT
	Casbin  *Casbin
	OpenAPI *OpenAPI
	MongoDB *MongoDB
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
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	PrivateKey         string
	Issuer             string
}

type Casbin struct {
	Model  string
	Policy string
}

type OpenAPI struct {
	Schema string
}

type MongoDB struct {
	URI                      string
	Username                 string
	Password                 string
	ReplicaSet               string
	ServerSelectionTimeoutMs time.Duration
	ConnectTimeoutMs         time.Duration
	SocketTimeoutMs          time.Duration // query timeout
}

// NewConfig creates a Config instance
func NewConfig() *Config {
	return &Config{
		Config:  libconfig.New("APP"),
		HTTP:    libhttp.DefaultConfig,
		Logging: liblog.DefaultConfig,
		BaseUrl: "http://localhost:1323",
		Admin: &Admin{
			Create:   false,
			Email:    "admin@example.com",
			Username: "admin",
			Password: "changeme",
		},
		OAuth2: &OAuth2{
			ClientId:     "",
			ClientSecret: "",
		},
		JWT: &JWT{
			AccessTokenExpiry:  10 * time.Minute,
			RefreshTokenExpiry: (30 * 24) * time.Hour,
			PrivateKey:         "./private-key.pem",
			Issuer:             "http://localhost:1323",
		},
		Casbin: &Casbin{
			Model:  "./casbin/model.conf",
			Policy: "./casbin/policy.csv",
		},
		OpenAPI: &OpenAPI{
			Schema: "./openapi/openapi.yaml",
		},
		MongoDB: &MongoDB{
			URI:                      "mongodb://localhost:27017",
			Username:                 "",
			Password:                 "",
			ReplicaSet:               "",
			ServerSelectionTimeoutMs: time.Millisecond * 5000,
			ConnectTimeoutMs:         time.Millisecond * 5000,
			SocketTimeoutMs:          time.Millisecond * 30000,
		},
	}
}

const (
	AppName = libconfig.AppName
	EnvName = libconfig.EnvName

	HTTPBindAddress = libhttp.HTTPBindAddress
	HTTPBindPort    = libhttp.HTTPBindPort
	HTTPLogRequests = libhttp.HTTPLogRequests

	BaseUrl = "base-url"

	AdminCreate   = "admin-create"
	AdminEmail    = "admin-email"
	AdminUsername = "admin-username"
	AdminPassword = "admin-password"

	OAuth2ClientId     = "oauth2-client-id"
	OAuth2ClientSecret = "oauth2-client-secret"

	JWTAccessTokenExpiry  = "jwt-access-token-expiry"
	JWTRefreshTokenExpiry = "jwt-refresh-token-expiry"
	JWTPrivateKey         = "jwt-private-key"
	JWTIssuer             = "jwt-issuer"

	CasbinModel  = "casbin-model"
	CasbinPolicy = "casbin-policy"

	OpenAPISchema = "openapi-schema"

	MongoDBURI                      = "mongodb-uri"
	MongoDBUsername                 = "mongodb-username"
	MongoDBPassword                 = "mongodb-password"
	MongoDBReplicaSet               = "mongodb-replica-set"
	MongoDBServerSelectionTimeoutMs = "mongodb-server-selection-timeout-ms"
	MongoDBConnectTimeoutMs         = "mongodb-connect-timeout-ms"
	MongoDBSocketTimeoutMs          = "mongodb-socket-timeout-ms"
)

// addFlags adds all the flags from the command line
func (c *Config) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.BaseUrl, BaseUrl, c.BaseUrl, "Base URL where the app will be served")

	fs.BoolVar(&c.Admin.Create, AdminCreate, c.Admin.Create, "Create admin")
	fs.StringVar(&c.Admin.Email, AdminEmail, c.Admin.Email, "Admin email")
	fs.StringVar(&c.Admin.Username, AdminUsername, c.Admin.Username, "Admin username")
	fs.StringVar(&c.Admin.Password, AdminPassword, c.Admin.Password, "Admin password")

	fs.StringVar(&c.OAuth2.ClientId, OAuth2ClientId, c.OAuth2.ClientId, "OAuth2 client id")
	fs.StringVar(&c.OAuth2.ClientSecret, OAuth2ClientSecret, c.OAuth2.ClientSecret, "OAuth2 client secret")

	fs.DurationVar(&c.JWT.AccessTokenExpiry, JWTAccessTokenExpiry, c.JWT.AccessTokenExpiry,
		"JWT access token expiry")
	fs.DurationVar(&c.JWT.RefreshTokenExpiry, JWTRefreshTokenExpiry, c.JWT.RefreshTokenExpiry,
		"JWT refresh token expiry")
	fs.StringVar(&c.JWT.PrivateKey, JWTPrivateKey, c.JWT.PrivateKey, "JWT private key file path")
	fs.StringVar(&c.JWT.Issuer, JWTIssuer, c.JWT.Issuer, "JWT issuer")

	fs.StringVar(&c.Casbin.Model, CasbinModel, c.Casbin.Model, "Casbin model file")
	fs.StringVar(&c.Casbin.Policy, CasbinPolicy, c.Casbin.Policy, "Casbin policy file")

	fs.StringVar(&c.OpenAPI.Schema, OpenAPISchema, c.OpenAPI.Schema, "OpenAPI schema file")

	fs.StringVar(&c.MongoDB.URI, MongoDBURI, c.MongoDB.URI, "MongoDB URI")
	fs.StringVar(&c.MongoDB.Username, MongoDBUsername, c.MongoDB.Username, "MongoDB username")
	fs.StringVar(&c.MongoDB.Password, MongoDBPassword, c.MongoDB.Password, "MongoDB password")
	fs.StringVar(&c.MongoDB.ReplicaSet, MongoDBReplicaSet, c.MongoDB.ReplicaSet, "MongoDB replica set")
	fs.DurationVar(&c.MongoDB.ServerSelectionTimeoutMs, MongoDBServerSelectionTimeoutMs,
		c.MongoDB.ServerSelectionTimeoutMs, "MongoDB server selection timeout ms")
	fs.DurationVar(&c.MongoDB.ConnectTimeoutMs, MongoDBConnectTimeoutMs, c.MongoDB.ConnectTimeoutMs,
		"MongoDB connect timeout ms")
	fs.DurationVar(&c.MongoDB.SocketTimeoutMs, MongoDBSocketTimeoutMs, c.MongoDB.SocketTimeoutMs,
		"MongoDB socket timeout ms")
}

func (c *Config) BindFlags() error {
	c.addFlags(pflag.CommandLine)
	c.Logging.BindFlags(pflag.CommandLine)
	c.HTTP.BindFlags(pflag.CommandLine)

	err := c.Config.BindFlags()
	if err != nil {
		return err
	}

	err = liblog.New(&liblog.Config{
		LogLevel:  viper.GetString(liblog.LogLevel),
		LogOutput: viper.GetString(liblog.LogOutput),
		LogWriter: viper.GetString(liblog.LogWriter),
	})
	if err != nil {
		return err
	}

	return nil
}
