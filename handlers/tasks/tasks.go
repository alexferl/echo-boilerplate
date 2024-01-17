package tasks

import (
	"context"
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

	body := &CreateTaskRequest{}
	if err := c.Bind(body); err != nil {
		logger.Error().Err(err).Msg("failed binding body")
		return err
	}

	token := c.Get("token").(jwt.Token)

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
		logger.Error().Err(err).Msg("failed to insert task")
		return err
	}

	pipeline := h.getPipeline(bson.D{{"_id", insert.InsertedID.(primitive.ObjectID)}}, 1, 0)
	result, err := h.Mapper.Aggregate(ctx, pipeline, Aggregates{})
	if err != nil {
		logger.Error().Err(err).Msg("failed getting task")
		return err
	}

	tasks := result.(Aggregates)
	if len(tasks) < 1 {
		logger.Error().Err(err).Msg("failed to retrieve inserted task")
		return err
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
	count, err := h.Mapper.Count(ctx, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed counting tasks")
		return err
	}

	pipeline := h.getPipeline(bson.D{{"deleted_at", bson.M{"$eq": nil}}}, limit, skip)
	result, err := h.Mapper.Aggregate(ctx, pipeline, Aggregates{})
	if err != nil {
		logger.Error().Err(err).Msg("failed getting tasks")
		return err
	}

	util.SetPaginationHeaders(c.Request(), c.Response().Header(), int(count), page, perPage)

	return h.Validate(c, http.StatusOK, &ListTasksResponse{Tasks: result.(Aggregates).Response()})
}
