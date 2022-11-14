package tasks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
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
	body := &UpdateTaskRequest{}
	if err := c.Bind(body); err != nil {
		return err
	}

	taskId := c.Param("id")
	token := c.Get("token").(jwt.Token)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	task, errResp := h.getTask(ctx, c, taskId, token)
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

	update, err := h.Mapper.UpdateById(ctx, taskId, task, []*TaskResponse{})
	if err != nil {
		return fmt.Errorf("failed updating task: %v", err)
	}

	res := update.([]*TaskResponse)
	if len(res) < 1 {
		return fmt.Errorf("failed to retrieve updated task: %v", err)
	}

	return h.Validate(c, http.StatusOK, res[0])
}

func (h *Handler) DeleteTask(c echo.Context) error {
	taskId := c.Param("id")
	token := c.Get("token").(jwt.Token)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	task, errResp := h.getTask(ctx, c, taskId, token)
	if errResp != nil {
		return errResp()
	}

	task.Delete(token.Subject())

	_, err := h.Mapper.UpdateById(ctx, taskId, task, nil)
	if err != nil {
		return fmt.Errorf("failed deleting task: %v", err)
	}

	return h.Validate(c, http.StatusNoContent, nil)
}
