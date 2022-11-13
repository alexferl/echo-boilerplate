package tasks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/handler"
	"github.com/alexferl/golib/http/router"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
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

func (h *Handler) getTask(ctx context.Context, c echo.Context, taskId string, token jwt.Token) (*Task, func() error) {
	result, err := h.Mapper.FindOneById(ctx, taskId, &Task{})
	if err != nil {
		if err == ErrTaskNotFound {
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

func (h *Handler) getAggregate(ctx context.Context, c echo.Context, filter any) (*ShortTask, func() error) {
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

func wrap(err error) func() error {
	return func() error { return err }
}
