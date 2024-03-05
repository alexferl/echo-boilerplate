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
	Mapper data.Mapper
}

func NewHandler(db *mongo.Client, openapi *openapi.Handler, mapper data.Mapper) handlers.IHandler {
	if mapper == nil {
		mapper = data.NewMapper(db, viper.GetString(config.AppName), "users")
	}

	if viper.GetBool(config.SuperUserCreate) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		filter := bson.D{{"username", viper.GetString(config.SuperUserUsername)}}
		_, err := mapper.FindOne(ctx, filter, &User{})
		if err != nil {
			if errors.Is(err, data.ErrNoDocuments) {
				log.Info().Msg("Creating super user")

				user := NewUserWithRole(viper.GetString(config.SuperUserEmail), viper.GetString(config.SuperUserUsername), SuperRole)
				user.Name = "Super User"
				user.Bio = "I am super."
				err = user.SetPassword(viper.GetString(config.SuperUserPassword))
				if err != nil {
					panic(fmt.Sprintf("failed setting super user password: %v", err))
				}

				user.Create(user.Id)

				ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				_, err = mapper.InsertOne(ctx, user, nil)
				if err != nil {
					panic(fmt.Sprintf("failed creating super user: %v", err))
				}
			} else {
				panic(fmt.Sprintf("failed getting super user: %v", err))
			}
		}
	}

	return &Handler{
		Handler: openapi,
		Mapper:  mapper,
	}
}

func (h *Handler) AddRoutes(s *server.Server) {
	s.Add(http.MethodGet, "/me", h.GetCurrentUser)
	s.Add(http.MethodPut, "/me", h.UpdateCurrentUser)
	s.Add(http.MethodPost, "/me/personal_access_tokens", h.CreatePersonalAccessToken)
	s.Add(http.MethodGet, "/me/personal_access_tokens", h.ListPersonalAccessTokens)
	s.Add(http.MethodGet, "/me/personal_access_tokens/:id", h.GetPersonalAccessToken)
	s.Add(http.MethodDelete, "/me/personal_access_tokens/:id", h.RevokePersonalAccessToken)
	s.Add(http.MethodGet, "/users/:id_or_username", h.GetUser)
	s.Add(http.MethodPut, "/users/:id_or_username", h.UpdateUser)
	s.Add(http.MethodPut, "/users/:id_or_username/status", h.UpdateUserStatus)
	s.Add(http.MethodGet, "/users", h.ListUsers)
}
