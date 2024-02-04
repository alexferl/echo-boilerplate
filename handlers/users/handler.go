package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers"
)

type Handler struct {
	*openapi.Handler
	Mapper data.IMapper
}

func NewHandler(db *mongo.Client, openapi *openapi.Handler, mapper data.IMapper) handlers.IHandler {
	if mapper == nil {
		mapper = data.NewMapper(db, viper.GetString(config.AppName), "users")
	}

	if viper.GetBool(config.AdminCreate) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		filter := bson.D{{"username", viper.GetString(config.AdminUsername)}}
		_, err := mapper.FindOne(ctx, filter, &User{})
		if err != nil {
			if errors.Is(err, data.ErrNoDocuments) {
				log.Info().Msg("Creating admin user")

				user := NewAdminUser(viper.GetString(config.AdminEmail), viper.GetString(config.AdminUsername))
				user.Name = "The Admin"
				user.Bio = "I am the admin!"
				err = user.SetPassword(viper.GetString(config.AdminPassword))
				if err != nil {
					panic(fmt.Sprintf("failed setting admin password: %v", err))
				}

				user.Create(user.Id)

				ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				_, err = mapper.InsertOne(ctx, user, nil)
				if err != nil {
					panic(fmt.Sprintf("failed creating admin user: %v", err))
				}
			} else {
				panic(fmt.Sprintf("failed getting admin user: %v", err))
			}
		}
	}

	return &Handler{
		Handler: openapi,
		Mapper:  mapper,
	}
}

func (h *Handler) AddRoutes(s *server.Server) {
	s.Add(http.MethodPost, "/auth/login", h.AuthLogIn)
	s.Add(http.MethodPost, "/auth/logout", h.AuthLogOut)
	s.Add(http.MethodPost, "/auth/refresh", h.AuthRefresh)
	s.Add(http.MethodPost, "/auth/signup", h.AuthSignUp)
	s.Add(http.MethodGet, "/auth/token", h.AuthToken)
	s.Add(http.MethodGet, "/oauth2/callback", h.OAuth2Callback)
	s.Add(http.MethodGet, "/oauth2/login", h.OAuth2LogIn)
	s.Add(http.MethodGet, "/me", h.GetMe)
	s.Add(http.MethodPut, "/me", h.UpdateMe)
	s.Add(http.MethodPost, "/me/personal_access_tokens", h.CreatePersonalAccessToken)
	s.Add(http.MethodGet, "/me/personal_access_tokens", h.ListPersonalAccessTokens)
	s.Add(http.MethodGet, "/me/personal_access_tokens/:id", h.GetPersonalAccessToken)
	s.Add(http.MethodDelete, "/me/personal_access_tokens/:id", h.RevokePersonalAccessToken)
	s.Add(http.MethodGet, "/users/:username", h.GetUsername)
	s.Add(http.MethodGet, "/users", h.ListUsers)
}
