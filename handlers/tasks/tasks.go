package tasks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/alexferl/echo-boilerplate/util"
)

type ShortTask struct {
	Id          string     `json:"id" bson:"id"`
	Title       string     `json:"title" bson:"title"`
	IsPrivate   bool       `json:"is_private" bson:"is_private"`
	IsCompleted bool       `json:"is_completed" bson:"is_completed"`
	CreatedAt   *time.Time `json:"created_at" bson:"created_at"`
	CreatedBy   *TaskUser  `json:"created_by" bson:"created_by"`
}

type CreateTaskPayload struct {
	Title     string `json:"title"`
	IsPrivate bool   `json:"is_private"`
}

func (h *Handler) CreateTask(c echo.Context) error {
	body := &CreateTaskPayload{}
	if err := c.Bind(body); err != nil {
		return fmt.Errorf("failed to bind: %v", err)
	}

	token := c.Get("token").(jwt.Token)

	newTask := NewTask()
	newTask.Create(token.Subject())
	newTask.Title = body.Title
	newTask.IsPrivate = body.IsPrivate

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := h.Mapper.Insert(ctx, newTask, []*ListTasks{})
	if err != nil {
		return fmt.Errorf("failed to insert task: %v", err)
	}

	tasks := result.([]*ListTasks)
	if len(tasks) < 1 {
		return fmt.Errorf("failed to retrieve inserted task: %v", err)
	}

	task := tasks[0]
	resp := ShortTask{
		Id:          task.Id,
		Title:       task.Title,
		IsPrivate:   task.IsPrivate,
		IsCompleted: task.IsCompleted,
		CreatedAt:   task.CreatedAt,
		CreatedBy:   task.CreatedBy,
	}

	return h.Validate(c, http.StatusOK, resp)
}

type ListTasks struct {
	Id          string     `json:"id" bson:"id"`
	Title       string     `json:"title" bson:"title"`
	IsPrivate   bool       `json:"is_private" bson:"is_private"`
	IsCompleted bool       `json:"is_completed" bson:"is_completed"`
	CompletedAt *time.Time `json:"completed_at" bson:"completed_at"`
	CompletedBy *TaskUser  `json:"completed_by" bson:"completed_by"`
	CreatedAt   *time.Time `json:"created_at" bson:"created_at"`
	CreatedBy   *TaskUser  `json:"created_by" bson:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy   *TaskUser  `json:"updated_by" bson:"updated_by"`
}

type ListTasksResp struct {
	Tasks []*ListTasks `json:"tasks"`
}

func (h *Handler) ListTasks(c echo.Context) error {
	page, perPage, limit, skip := util.ParsePaginationParams(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := h.Mapper.Count(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed counting tasks: %v", err)
	}

	result, err := h.Mapper.Aggregate(ctx, nil, limit, skip, []*ListTasks{})
	if err != nil {
		return fmt.Errorf("failed getting tasks: %v", err)
	}

	uri := fmt.Sprintf("http://%s%s", c.Request().Host, c.Request().URL.Path)
	util.SetPaginationHeaders(c.Response().Header(), int(count), page, perPage, uri)

	return h.Validate(c, http.StatusOK, &ListTasksResp{Tasks: result.([]*ListTasks)})
}
