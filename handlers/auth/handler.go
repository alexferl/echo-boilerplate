package auth

import (
	"net/http"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/spf13/viper"
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

	return &Handler{
		Handler: openapi,
		Mapper:  mapper,
	}
}

func (h *Handler) AddRoutes(s *server.Server) {
	s.Add(http.MethodPost, "/auth/login", h.Login)
	s.Add(http.MethodPost, "/auth/logout", h.Logout)
	s.Add(http.MethodPost, "/auth/refresh", h.AuthRefresh)
	s.Add(http.MethodPost, "/auth/signup", h.AuthSignUp)
	s.Add(http.MethodGet, "/auth/token", h.AuthToken)
	s.Add(http.MethodGet, "/oauth2/callback", h.OAuth2Callback)
	s.Add(http.MethodGet, "/oauth2/login", h.OAuth2LogIn)
}
