package tasks

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
)

type Handler struct {
	*openapi.Handler
	Mapper data.Mapper
	Db     *mongo.Client
}

func NewHandler(db *mongo.Client, openapi *openapi.Handler, mapper data.Mapper) handlers.BaseHandler {
	if mapper == nil {
		mapper = NewMapper(db)
	}

	return &Handler{
		Handler: openapi,
		Mapper:  mapper,
		Db:      db,
	}
}

func (h *Handler) AddRoutes(s *server.Server) {
	s.Add(http.MethodPost, "/tasks", h.CreateTask)
	s.Add(http.MethodGet, "/tasks", h.ListTasks)
	s.Add(http.MethodGet, "/tasks/:id", h.GetTask)
	s.Add(http.MethodPatch, "/tasks/:id", h.UpdateTask)
	s.Add(http.MethodDelete, "/tasks/:id", h.DeleteTask)
}

func (h *Handler) getTask(ctx context.Context, c echo.Context, taskId string, token jwt.Token) (*Task, func() error) {
	result, err := h.Mapper.FindOneById(ctx, taskId, &Task{})
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			return nil, wrap(h.Validate(c, http.StatusNotFound, echo.Map{"message": "task not found"}))
		}
		return nil, wrap(fmt.Errorf("failed getting task: %v", err))
	}

	task := result.(*Task)
	if task.DeletedAt != nil {
		return nil, wrap(h.Validate(c, http.StatusGone, echo.Map{"message": "task was deleted"}))
	}

	if token != nil {
		if token.Subject() != task.CreatedBy && !util.HasRole(token, users.AdminRole.String()) {
			return nil, wrap(h.Validate(c, http.StatusForbidden, echo.Map{"message": "you don't have access"}))
		}
	}

	return task, nil
}

func (h *Handler) getAggregate(ctx context.Context, c echo.Context) (*TaskResponse, func() error) {
	filter := bson.D{{"id", c.Param("id")}}
	result, err := h.Mapper.Aggregate(ctx, filter, 1, 0, []*TaskResponse{})
	if err != nil {
		return nil, wrap(fmt.Errorf("failed getting task: %v", err))
	}

	tasks := result.([]*TaskResponse)
	if len(tasks) < 1 {
		return nil, wrap(h.Validate(c, http.StatusNotFound, echo.Map{"message": "task not found"}))
	}

	task := tasks[0]
	if task.DeletedAt != nil {
		return nil, wrap(h.Validate(c, http.StatusGone, echo.Map{"message": "task was deleted"}))
	}

	return task, nil
}

func wrap(err error) func() error {
	return func() error { return err }
}
