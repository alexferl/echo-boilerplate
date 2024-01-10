package tasks

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexferl/echo-boilerplate/util"
)

type CreateTaskRequest struct {
	Title string `json:"title"`
}

func (h *Handler) CreateTask(c echo.Context) error {
	body := &CreateTaskRequest{}
	if err := c.Bind(body); err != nil {
		return fmt.Errorf("failed to bind: %v", err)
	}

	token := c.Get("token").(jwt.Token)

	newTask := NewTask()
	newTask.Create(token.Subject())
	newTask.Title = body.Title

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insert, err := h.Mapper.InsertOne(ctx, newTask)
	if err != nil {
		return fmt.Errorf("failed to insert task: %v", err)
	}

	pipeline := h.getPipeline(bson.D{{"_id", insert.InsertedID.(primitive.ObjectID)}}, 1, 0)
	result, err := h.Mapper.Aggregate(ctx, pipeline, []*TaskResponse{})
	if err != nil {
		return fmt.Errorf("failed getting tasks: %v", err)
	}

	tasks := result.([]*TaskResponse)
	if len(tasks) < 1 {
		return fmt.Errorf("failed to retrieve inserted task: %v", err)
	}

	return h.Validate(c, http.StatusOK, tasks[0])
}

type ListTasksResponse struct {
	Tasks []*TaskResponse `json:"tasks"`
}

func (h *Handler) ListTasks(c echo.Context) error {
	page, perPage, limit, skip := util.ParsePaginationParams(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := h.Mapper.Count(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed counting tasks: %v", err)
	}

	pipeline := h.getPipeline(bson.D{{"deleted_at", bson.M{"$eq": nil}}}, limit, skip)
	result, err := h.Mapper.Aggregate(ctx, pipeline, []*TaskResponse{})
	if err != nil {
		return fmt.Errorf("failed getting tasks: %v", err)
	}

	uri := fmt.Sprintf("http://%s%s", c.Request().Host, c.Request().URL.Path)
	util.SetPaginationHeaders(c.Response().Header(), int(count), page, perPage, uri)

	return h.Validate(c, http.StatusOK, &ListTasksResponse{Tasks: result.([]*TaskResponse)})
}
