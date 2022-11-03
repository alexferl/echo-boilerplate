package users_test

import (
	"testing"
	"time"

	openapi "github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/server"
	"go.mongodb.org/mongo-driver/mongo"

	app "github.com/alexferl/echo-boilerplate"
	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/mocks"
)

var overrides = map[string]any{
	config.JWTPrivateKey:         "../../private-key.pem",
	config.CasbinModel:           "../../casbin/model.conf",
	config.CasbinPolicy:          "../../casbin/policy.csv",
	config.OpenAPISchema:         "../../openapi/openapi.yaml",
	config.JWTAccessTokenExpiry:  10 * time.Minute,
	config.JWTRefreshTokenExpiry: (30 * 24) * time.Hour,
	config.HTTPLogRequests:       false,
}

func getMapperAndServer(t *testing.T) (*mocks.Mapper, *server.Server) {
	mapper := mocks.NewMapper(t)
	h := users.NewHandler(&mongo.Client{}, openapi.NewHandler(), mapper)
	s := app.NewServerWithOverrides(overrides, h)
	return mapper, s
}
