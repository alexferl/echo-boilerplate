package tasks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
)

func (h *Handler) GetTask(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"id", c.Param("id")}}
	t, errResp := h.getAggregate(ctx, c, filter)
	if errResp != nil {
		return errResp()
	}

	return h.Validate(c, http.StatusOK, t)
}

type TaskPatch struct {
	Title       string `json:"title" bson:"title"`
	IsPrivate   bool   `json:"is_private" bson:"is_private"`
	IsCompleted bool   `json:"is_completed" bson:"is_completed"`
}

func (h *Handler) UpdateTask(c echo.Context) error {
	body := &TaskPatch{}
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

	if body.IsPrivate != task.IsPrivate {
		task.IsPrivate = body.IsPrivate
	}

	if body.IsCompleted != task.IsCompleted {
		task.IsCompleted = body.IsCompleted
	}

	task.Update(token.Subject())

	update, err := h.Mapper.UpdateById(ctx, taskId, task, []*TaskWithUsers{})
	if err != nil {
		return fmt.Errorf("failed updating task: %v", err)
	}

	res := update.([]*TaskWithUsers)
	if len(res) > 0 {
		updated := res[0]
		t := ShortTask{
			Id:          updated.Id,
			Title:       updated.Title,
			IsPrivate:   updated.IsPrivate,
			IsCompleted: updated.IsCompleted,
			CreatedAt:   updated.CreatedAt,
			CreatedBy:   updated.CreatedBy,
		}

		return h.Validate(c, http.StatusOK, t)
	}

	return fmt.Errorf("fail")
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

func wrap(err error) func() error {
	return func() error { return err }
}
