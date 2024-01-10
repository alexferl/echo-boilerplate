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
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
)

type Handler struct {
	*openapi.Handler
	Mapper data.IMapper
	Db     *mongo.Client
}

func NewHandler(db *mongo.Client, openapi *openapi.Handler, mapper data.IMapper) handlers.IHandler {
	if mapper == nil {
		mapper = data.NewMapper(db, viper.GetString(config.AppName), "tasks")
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
	s.Add(http.MethodPut, "/tasks/:id", h.UpdateTask)
	s.Add(http.MethodDelete, "/tasks/:id", h.DeleteTask)
}

func (h *Handler) getTask(ctx context.Context, c echo.Context, taskId string, token jwt.Token) (*Task, func() error) {
	result, err := h.Mapper.FindOneById(ctx, taskId, &Task{})
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return nil, util.Wrap(h.Validate(c, http.StatusNotFound, echo.Map{"message": "task not found"}))
		}
		return nil, util.Wrap(fmt.Errorf("failed getting task: %v", err))
	}

	task := result.(*Task)
	if task.DeletedAt != nil {
		return nil, util.Wrap(h.Validate(c, http.StatusGone, echo.Map{"message": "task was deleted"}))
	}

	if token != nil {
		if token.Subject() != task.CreatedBy && !util.HasRole(token, users.AdminRole.String()) {
			return nil, util.Wrap(h.Validate(c, http.StatusForbidden, echo.Map{"message": "you don't have access"}))
		}
	}

	return task, nil
}

func (h *Handler) getPipeline(filter any, limit int, skip int) mongo.Pipeline {
	if filter == nil {
		filter = bson.D{}
	}

	return mongo.Pipeline{
		{{"$match", filter}},
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "created_by",
			"foreignField": "id",
			"as":           "created_by",
		}}},
		{{"$unwind", "$created_by"}},
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "updated_by",
			"foreignField": "id",
			"as":           "updated_by",
		}}},
		{{
			"$unwind", bson.D{
				{"path", "$updated_by"},
				{"preserveNullAndEmptyArrays", true},
			},
		}},
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "completed_by",
			"foreignField": "id",
			"as":           "completed_by",
		}}},
		{{
			"$unwind", bson.D{
				{"path", "$completed_by"},
				{"preserveNullAndEmptyArrays", true},
			},
		}},
		{{"$limit", skip + limit}},
		{{"$skip", skip}},
	}
}

func (h *Handler) getAggregate(ctx context.Context, c echo.Context) (*TaskResponse, func() error) {
	pipeline := h.getPipeline(bson.D{{"id", c.Param("id")}}, 1, 0)
	result, err := h.Mapper.Aggregate(ctx, pipeline, []*TaskResponse{})
	if err != nil {
		return nil, util.Wrap(fmt.Errorf("failed getting task: %v", err))
	}

	tasks := result.([]*TaskResponse)
	if len(tasks) < 1 {
		return nil, util.Wrap(h.Validate(c, http.StatusNotFound, echo.Map{"message": "task not found"}))
	}

	task := tasks[0]
	if task.DeletedAt != nil {
		return nil, util.Wrap(h.Validate(c, http.StatusGone, echo.Map{"message": "task was deleted"}))
	}

	return task, nil
}
