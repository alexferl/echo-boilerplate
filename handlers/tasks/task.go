package tasks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *Handler) GetTask(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	task, errResp := h.getAggregate(ctx, c)
	if errResp != nil {
		return errResp()
	}

	return h.Validate(c, http.StatusOK, task)
}

type UpdateTaskRequest struct {
	Title     string `json:"title" bson:"title"`
	Completed bool   `json:"completed" bson:"completed"`
}

func (h *Handler) UpdateTask(c echo.Context) error {
	id := c.Param("id")
	logger := c.Get("logger").(zerolog.Logger)
	token := c.Get("token").(jwt.Token)

	body := &UpdateTaskRequest{}
	if err := c.Bind(body); err != nil {
		logger.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	task, errResp := h.getTask(ctx, c, id, token)
	if errResp != nil {
		return errResp()
	}

	if body.Title != "" {
		task.Title = body.Title
	}

	if body.Completed != task.Completed {
		if body.Completed {
			task.Complete(token.Subject())
		} else {
			task.Incomplete()
		}
	}

	task.Update(token.Subject())

	_, err := h.Mapper.UpdateOneById(ctx, id, task)
	if err != nil {
		logger.Error().Err(err).Msg("failed updating task")
		return err
	}

	pipeline := h.getPipeline(bson.D{{"id", id}}, 1, 0)
	res, err := h.Mapper.Aggregate(ctx, pipeline, Aggregates{})
	if err != nil {
		logger.Error().Err(err).Msg("failed getting task")
		return err
	}

	resp := res.(Aggregates)
	if len(resp) < 1 {
		msg := "failed retrieving updated task"
		logger.Error().Msg(msg)
		return fmt.Errorf(msg)
	}

	return h.Validate(c, http.StatusOK, resp[0].Response())
}

func (h *Handler) DeleteTask(c echo.Context) error {
	id := c.Param("id")
	logger := c.Get("logger").(zerolog.Logger)
	token := c.Get("token").(jwt.Token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	task, errResp := h.getTask(ctx, c, id, token)
	if errResp != nil {
		return errResp()
	}

	task.Delete(token.Subject())

	_, err := h.Mapper.UpdateOneById(ctx, id, task, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed deleting task")
		return err
	}

	return h.Validate(c, http.StatusNoContent, nil)
}
