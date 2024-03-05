package tasks

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"

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

	model := h.model.New()
	task, err := model.Create(ctx, token.Subject(), *body)
	if err != nil {
		logger.Error().Err(err).Msg("failed creating task")
		return err
	}

	return h.Validate(c, http.StatusOK, task.Response())
}

type ListTasksResponse struct {
	Tasks []*Response `json:"tasks"`
}

func (h *Handler) ListTasks(c echo.Context) error {
	logger := c.Get("logger").(zerolog.Logger)
	page, perPage, limit, skip := util.ParsePaginationParams(c)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	completed := c.QueryParams()["completed"]
	if len(completed) > 0 {
		arr := bson.A{}
		for _, i := range completed {
			s := strings.ToLower(i)
			if s == "true" {
				arr = append(arr, true)
			} else if s == "false" {
				arr = append(arr, false)
			}
		}
		filter["completed"] = bson.M{"$in": arr}
	}
	createdBy := c.QueryParam("created_by")
	if createdBy != "" {
		filter["created_by"] = createdBy
	}
	query := c.QueryParams()["q"]
	if len(query) > 0 {
		filter["$text"] = bson.M{"$search": strings.Join(query, " ")}
	}

	count, tasks, err := h.model.Find(ctx, filter, limit, skip)
	if err != nil {
		logger.Error().Err(err).Msg("failed finding tasks")
		return err
	}

	util.SetPaginationHeaders(c.Request(), c.Response().Header(), int(count), page, perPage)

	return h.Validate(c, http.StatusOK, &ListTasksResponse{Tasks: tasks.Response()})
}
