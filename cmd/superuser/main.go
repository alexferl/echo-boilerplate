package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/alexferl/golib/config"
	"github.com/alexferl/golib/database/mongodb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/mappers"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/services"
)

type Config struct {
	Config  *config.Config
	MongoDB *mongodb.Config
	Super   *Super
}

type Super struct {
	Email    string
	Name     string
	Password string
	Username string
}

func New() *Config {
	return &Config{
		Config:  config.New("APP"),
		MongoDB: mongodb.DefaultConfig,
		Super: &Super{
			Email:    "super@example.com",
			Name:     "Super",
			Password: "",
			Username: "super",
		},
	}
}

const (
	SuperEmail    = "email"
	SuperName     = "name"
	SuperPassword = "password"
	SuperUsername = "username"
)

func (c *Config) addFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Super.Email, SuperEmail, c.Super.Email, "Superuser email")
	fs.StringVar(&c.Super.Name, SuperName, c.Super.Name, "Superuser display name")
	fs.StringVar(&c.Super.Password, SuperPassword, c.Super.Password, "Superuser password")
	fs.StringVar(&c.Super.Username, SuperUsername, c.Super.Username, "Superuser name")
}

func (c *Config) BindFlags() {
	c.addFlags(pflag.CommandLine)
	c.MongoDB.BindFlags(pflag.CommandLine)

	err := c.Config.BindFlags()
	if err != nil {
		log.Fatal().Err(err).Msg("failed binding flags")
	}

	if viper.GetString(SuperPassword) == "" {
		log.Fatal().Msg("password is unset!")
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	c := New()
	c.BindFlags()

	client, err := data.MewMongoClient()
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating mongo client")
	}

	mapper := mappers.NewUser(client)
	svc := services.NewUser(mapper)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := svc.FindOneByEmailOrUsername(ctx, viper.GetString(SuperEmail), viper.GetString(SuperUsername))
	if res != nil {
		log.Fatal().Msg("username or email already in-use")
	}
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			email := viper.GetString(SuperEmail)
			name := viper.GetString(SuperName)
			username := viper.GetString(SuperUsername)

			log.Info().
				Str("name", name).
				Str("username", username).
				Str("email", email).
				Msg("creating superuser")

			user := models.NewUserWithRole(email, username, models.SuperRole)
			user.Name = name

			err = user.SetPassword(viper.GetString(SuperPassword))
			if err != nil {
				log.Fatal().Err(err).Msg("failed setting superuser password")
			}

			_, err = svc.Create(ctx, user)
			if err != nil {
				log.Fatal().Err(err).Msg("failed creating superuser")
			}

			log.Info().Msg("done")
		} else {
			log.Fatal().Err(err).Msg("failed getting superuser")
		}
	}
}
