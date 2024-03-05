package tasks

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
	model Repository
}

func NewHandler(client *mongo.Client, openapi *openapi.Handler, modeler Repository) handlers.IHandler {
	if modeler == nil {
		modeler = NewModel(data.NewMapper(client, viper.GetString(config.AppName), "tasks"))
	}

	return &Handler{
		Handler: openapi,
		model:   modeler,
	}
}

func (h *Handler) AddRoutes(s *server.Server) {
	s.Add(http.MethodPost, "/tasks", h.CreateTask)
	s.Add(http.MethodGet, "/tasks", h.ListTasks)
	s.Add(http.MethodGet, "/tasks/:id", h.GetTask)
	s.Add(http.MethodPut, "/tasks/:id", h.UpdateTask)
	s.Add(http.MethodPut, "/tasks/:id/transition", h.TransitionTask)
	s.Add(http.MethodDelete, "/tasks/:id", h.DeleteTask)
}
