package tasks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/handler"
	"github.com/alexferl/golib/http/router"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexferl/echo-boilerplate/data"
)

type Handler struct {
	*openapi.Handler
	Mapper data.Mapper
	Db     *mongo.Client
}

func NewHandler(db *mongo.Client, openapi *openapi.Handler, mapper data.Mapper) handler.Handler {
	if mapper == nil {
		mapper = NewMapper(db)
	}

	return &Handler{
		Handler: openapi,
		Mapper:  mapper,
		Db:      db,
	}
}

func (h *Handler) GetRoutes() []*router.Route {
	return []*router.Route{
		{Name: "CreateTask", Method: http.MethodPost, Pattern: "/tasks", HandlerFunc: h.CreateTask},
		{Name: "ListTasks", Method: http.MethodGet, Pattern: "/tasks", HandlerFunc: h.ListTasks},
		{Name: "GetTask", Method: http.MethodGet, Pattern: "/tasks/:id", HandlerFunc: h.GetTask},
		{Name: "UpdateTask", Method: http.MethodPatch, Pattern: "/tasks/:id", HandlerFunc: h.UpdateTask},
		{Name: "DeleteTask", Method: http.MethodDelete, Pattern: "/tasks/:id", HandlerFunc: h.DeleteTask},
	}
}

func (h *Handler) getAggregate(c echo.Context, filter any) (*ShortTask, func() error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := h.Mapper.Aggregate(ctx, filter, 1, 0, []*TaskWithUsers{})
	if err != nil {
		return nil, wrap(fmt.Errorf("failed getting task: %v", err))
	}

	tasks := result.([]*TaskWithUsers)
	if len(tasks) < 1 {
		return nil, wrap(h.Validate(c, http.StatusNotFound, echo.Map{"message": "task not found"}))
	}

	task := tasks[0]
	if task.DeletedAt != nil {
		return nil, wrap(h.Validate(c, http.StatusGone, echo.Map{"message": "task was deleted"}))
	}

	t := &ShortTask{
		Id:          task.Id,
		Title:       task.Title,
		IsPrivate:   task.IsPrivate,
		IsCompleted: task.IsCompleted,
		CreatedAt:   task.CreatedAt,
		CreatedBy:   task.CreatedBy,
	}

	return t, nil
}
