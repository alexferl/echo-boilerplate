package tasks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
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
	body := &UpdateTaskRequest{}
	if err := c.Bind(body); err != nil {
		return err
	}

	id := c.Param("id")
	token := c.Get("token").(jwt.Token)
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
		return fmt.Errorf("failed updating task: %v", err)
	}

	pipeline := h.getPipeline(bson.D{{"id", id}}, 1, 0)
	result, err := h.Mapper.Aggregate(ctx, pipeline, []*TaskResponse{})
	if err != nil {
		return fmt.Errorf("failed getting tasks: %v", err)
	}

	res := result.([]*TaskResponse)
	if len(res) < 1 {
		return fmt.Errorf("failed to retrieve updated task: %v", err)
	}

	return h.Validate(c, http.StatusOK, res[0])
}

func (h *Handler) DeleteTask(c echo.Context) error {
	id := c.Param("id")
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
		return fmt.Errorf("failed deleting task: %v", err)
	}

	return h.Validate(c, http.StatusNoContent, nil)
}
