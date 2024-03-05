package tasks

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog"
)

func (h *Handler) GetTask(c echo.Context) error {
	id := c.Param("id")
	token := c.Get("token").(jwt.Token)
	logger := c.Get("logger").(zerolog.Logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	model, err := h.model.Load(ctx, id, token)
	if err != nil {
		logger.Error().Err(err).Msg("failed loading task")
		return err
	}

	task, err := model.Read(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed reading task")
		return err
	}

	return h.Validate(c, http.StatusOK, task.Response())
}

type UpdateTaskRequest struct {
	Title *string `json:"title"`
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

	model, err := h.model.Load(ctx, id, token)
	if err != nil {
		logger.Error().Err(err).Msg("failed loading task")
		return err
	}

	if body.Title != nil {
		model.Title = *body.Title
	}

	task, err := model.Update(ctx, token.Subject())
	if err != nil {
		logger.Error().Err(err).Msg("failed updating task")
		return err
	}

	return h.Validate(c, http.StatusOK, task.Response())
}

type TransitionTaskRequest struct {
	Completed bool `json:"completed"`
}

func (h *Handler) TransitionTask(c echo.Context) error {
	id := c.Param("id")
	logger := c.Get("logger").(zerolog.Logger)
	token := c.Get("token").(jwt.Token)

	body := &TransitionTaskRequest{}
	if err := c.Bind(body); err != nil {
		logger.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	model, err := h.model.Load(ctx, id, token)
	if err != nil {
		logger.Error().Err(err).Msg("failed loading task")
		return err
	}

	if body.Completed != model.Completed {
		if body.Completed {
			model.Complete(token.Subject())
		} else {
			model.Incomplete()
		}
	}

	task, err := model.Update(ctx, token.Subject())
	if err != nil {
		logger.Error().Err(err).Msg("failed updating task")
		return err
	}

	return h.Validate(c, http.StatusOK, task.Response())
}

func (h *Handler) DeleteTask(c echo.Context) error {
	id := c.Param("id")
	logger := c.Get("logger").(zerolog.Logger)
	token := c.Get("token").(jwt.Token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	model, err := h.model.Load(ctx, id, token)
	if err != nil {
		logger.Error().Err(err).Msg("failed loading task")
		return err
	}

	err = model.Delete(ctx, token.Subject())
	if err != nil {
		logger.Error().Err(err).Msg("failed deleting task")
		return err
	}

	return h.Validate(c, http.StatusNoContent, nil)
}
