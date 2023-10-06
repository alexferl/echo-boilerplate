package tasks_test

import (
	"testing"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"go.mongodb.org/mongo-driver/mongo"

	app "github.com/alexferl/echo-boilerplate"
	"github.com/alexferl/echo-boilerplate/handlers/tasks"
	"github.com/alexferl/echo-boilerplate/mocks"
	_ "github.com/alexferl/echo-boilerplate/testing"
)

func getMapperAndServer(t *testing.T) (*mocks.Mapper, *server.Server) {
	mapper := mocks.NewMapper(t)
	h := tasks.NewHandler(&mongo.Client{}, openapi.NewHandler(), mapper)
	s := app.NewTestServer(h)
	return mapper, s
}
