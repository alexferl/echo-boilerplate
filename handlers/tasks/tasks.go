package tasks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/alexferl/echo-boilerplate/util"
)

type CreateTaskRequest struct {
	Title string `json:"title"`
}

func (h *Handler) CreateTask(c echo.Context) error {
	logger := c.Get("logger").(zerolog.Logger)
	token := c.Get("token").(jwt.Token)

	body := &CreateTaskRequest{}
	if err := c.Bind(body); err != nil {
		logger.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	seq, err := h.Mapper.GetNextSequence(ctx, "tasks")
	if err != nil {
		logger.Error().Err(err).Msg("failed getting next sequence")
		return err
	}

	newTask := NewTask(seq.String())
	newTask.Create(token.Subject())
	newTask.Title = body.Title

	insert, err := h.Mapper.InsertOne(ctx, newTask)
	if err != nil {
		logger.Error().Err(err).Msg("failed inserting task")
		return err
	}

	pipeline := h.getPipeline(bson.D{{"_id", insert.InsertedID.(primitive.ObjectID)}}, 1, 0)
	res, err := h.Mapper.Aggregate(ctx, pipeline, Aggregates{})
	if err != nil {
		logger.Error().Err(err).Msg("failed getting task")
		return err
	}

	tasks := res.(Aggregates)
	if len(tasks) < 1 {
		msg := "failed retrieving inserted task"
		logger.Error().Msg(msg)
		return fmt.Errorf(msg)
	}

	return h.Validate(c, http.StatusOK, tasks[0].Response())
}

type ListTasksResponse struct {
	Tasks []*Response `json:"tasks"`
}

func (h *Handler) ListTasks(c echo.Context) error {
	logger := c.Get("logger").(zerolog.Logger)
	page, perPage, limit, skip := util.ParsePaginationParams(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pipeline := h.getPipeline(bson.D{{"deleted_at", bson.M{"$eq": nil}}}, limit, skip)
	res, err := h.Mapper.Aggregate(ctx, pipeline, Aggregates{})
	if err != nil {
		logger.Error().Err(err).Msg("failed getting tasks")
		return err
	}

	tasks := res.(Aggregates)
	util.SetPaginationHeaders(c.Request(), c.Response().Header(), len(tasks), page, perPage)

	return h.Validate(c, http.StatusOK, &ListTasksResponse{Tasks: tasks.Response()})
}
