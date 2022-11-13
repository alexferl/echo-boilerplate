package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/handler"
	"github.com/alexferl/golib/http/router"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
)

type Handler struct {
	*openapi.Handler
	Mapper data.Mapper
}

func NewHandler(db *mongo.Client, openapi *openapi.Handler, mapper data.Mapper) handler.Handler {
	if mapper == nil {
		mapper = NewMapper(db)
	}

	if viper.GetBool(config.AdminCreate) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		filter := bson.D{{"username", viper.GetString(config.AdminUsername)}}
		_, err := mapper.FindOne(ctx, filter, &User{})
		if err != nil {
			if err == ErrUserNotFound {
				log.Info().Msg("Creating admin user")

				user := NewAdminUser(viper.GetString(config.AdminEmail), viper.GetString(config.AdminUsername))
				err = user.SetPassword(viper.GetString(config.AdminPassword))
				user.Create(user.Id)
				if err != nil {
					panic(fmt.Sprintf("failed setting admin password: %v", err))
				}

				ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				_, err = mapper.Insert(ctx, user, nil)
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

func (h *Handler) GetRoutes() []*router.Route {
	return []*router.Route{
		{Name: "AuthSignup", Method: http.MethodPost, Pattern: "/auth/signup", HandlerFunc: h.AuthSignUp},
		{Name: "AuthLogin", Method: http.MethodPost, Pattern: "/auth/login", HandlerFunc: h.AuthLogin},
		{Name: "AuthRefresh", Method: http.MethodPost, Pattern: "/auth/refresh", HandlerFunc: h.AuthRefresh},
		{Name: "AuthLogout", Method: http.MethodPost, Pattern: "/auth/logout", HandlerFunc: h.AuthLogout},
		{Name: "OAuth2Login", Method: http.MethodGet, Pattern: "/oauth2/login", HandlerFunc: h.OAuth2Login},
		{Name: "OAuth2Callback", Method: http.MethodGet, Pattern: "/oauth2/callback", HandlerFunc: h.OAuth2Callback},
		{Name: "GetUser", Method: http.MethodGet, Pattern: "/user", HandlerFunc: h.GetUser},
		{Name: "UpdateUser", Method: http.MethodPatch, Pattern: "/user", HandlerFunc: h.UpdateUser},
		{Name: "GetUsername", Method: http.MethodGet, Pattern: "/users/:username", HandlerFunc: h.GetUsername},
		{Name: "ListUsers", Method: http.MethodGet, Pattern: "/users", HandlerFunc: h.ListUsers},
	}
}