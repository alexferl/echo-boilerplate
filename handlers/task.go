package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog/log"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/util"
)

type TaskService interface {
	Create(ctx context.Context, id string, data *models.Task) (*models.Task, error)
	Read(ctx context.Context, id string) (*models.Task, error)
	Update(ctx context.Context, id string, data *models.Task) (*models.Task, error)
	Delete(ctx context.Context, id string, data *models.Task) error

	Find(ctx context.Context, params *models.TaskSearchParams) (int64, models.Tasks, error)
}

type TaskHandler struct {
	*openapi.Handler
	svc TaskService
}

func NewTaskHandler(openapi *openapi.Handler, svc TaskService) *TaskHandler {
	return &TaskHandler{
		Handler: openapi,
		svc:     svc,
	}
}

func (h *TaskHandler) Register(s *server.Server) {
	s.Add(http.MethodPost, "/tasks", h.create)
	s.Add(http.MethodGet, "/tasks", h.list)
	s.Add(http.MethodGet, "/tasks/:id", h.get)
	s.Add(http.MethodPut, "/tasks/:id", h.update)
	s.Add(http.MethodPut, "/tasks/:id/transition", h.transition)
	s.Add(http.MethodDelete, "/tasks/:id", h.delete)
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}

func (h *TaskHandler) create(c echo.Context) error {
	token := c.Get("token").(jwt.Token)

	body := &CreateTaskRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	model := models.NewTask()
	model.Title = body.Title

	task, err := h.svc.Create(ctx, token.Subject(), model)
	if err != nil {
		log.Error().Err(err).Msg("failed creating task")
		return err
	}

	return h.Validate(c, http.StatusOK, task.Response())
}

type ListTasksResponse struct {
	Tasks []models.TaskResponse `json:"tasks"`
}

func (h *TaskHandler) list(c echo.Context) error {
	page, perPage, limit, skip := util.ParsePaginationParams(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	params := &models.TaskSearchParams{
		Completed: c.QueryParams()["completed"],
		CreatedBy: c.QueryParam("created_by"),
		Queries:   c.QueryParams()["q"],
		Limit:     limit,
		Skip:      skip,
	}

	count, tasks, err := h.svc.Find(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed getting tasks")
		return err
	}

	util.SetPaginationHeaders(c.Request(), c.Response().Header(), int(count), page, perPage)

	return h.Validate(c, http.StatusOK, &ListTasksResponse{Tasks: tasks.Response()})
}

func (h *TaskHandler) get(c echo.Context) error {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	task, err := h.svc.Read(ctx, id)
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return h.Validate(c, http.StatusNotFound, echo.Map{"message": "task not found"})
		}
		log.Error().Err(err).Msg("failed getting user")
		return err
	}

	if task.DeletedAt != nil {
		return h.Validate(c, http.StatusGone, echo.Map{"message": "task was deleted"})
	}

	return h.Validate(c, http.StatusOK, task.Response())
}

type UpdateTaskRequest struct {
	Title *string `json:"title"`
}

func (h *TaskHandler) update(c echo.Context) error {
	id := c.Param("id")
	token := c.Get("token").(jwt.Token)

	body := &UpdateTaskRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	task, err := h.svc.Read(ctx, id)
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return h.Validate(c, http.StatusNotFound, echo.Map{"message": "task not found"})
		}
		log.Error().Err(err).Msg("failed getting task")
		return err
	}

	if task.DeletedAt != nil {
		return h.Validate(c, http.StatusGone, echo.Map{"message": "task was deleted"})
	}

	if token.Subject() != task.CreatedBy.(*models.User).Id && !util.HasRoles(token, models.AdminRole.String(), models.SuperRole.String()) {
		return h.Validate(c, http.StatusForbidden, echo.Map{"message": "you don't have access"})
	}

	if body.Title != nil {
		task.Title = *body.Title
	}

	res, err := h.svc.Update(ctx, token.Subject(), task)
	if err != nil {
		log.Error().Err(err).Msg("failed updating task")
		return err
	}

	return h.Validate(c, http.StatusOK, res.Response())
}

type TransitionTaskRequest struct {
	Completed *bool `json:"completed"`
}

func (h *TaskHandler) transition(c echo.Context) error {
	id := c.Param("id")
	token := c.Get("token").(jwt.Token)

	body := &TransitionTaskRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	task, err := h.svc.Read(ctx, id)
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return h.Validate(c, http.StatusNotFound, echo.Map{"message": "task not found"})
		}
		log.Error().Err(err).Msg("failed getting task")
		return err
	}

	if task.DeletedAt != nil {
		return h.Validate(c, http.StatusGone, echo.Map{"message": "task was deleted"})
	}

	if *body.Completed != task.Completed {
		if *body.Completed {
			task.Complete(token.Subject())
		} else {
			task.Incomplete()
		}
	}

	res, err := h.svc.Update(ctx, token.Subject(), task)
	if err != nil {
		log.Error().Err(err).Msg("failed updating task")
		return err
	}

	return h.Validate(c, http.StatusOK, res.Response())
}

func (h *TaskHandler) delete(c echo.Context) error {
	id := c.Param("id")
	token := c.Get("token").(jwt.Token)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	task, err := h.svc.Read(ctx, id)
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return h.Validate(c, http.StatusNotFound, echo.Map{"message": "task not found"})
		}
		log.Error().Err(err).Msg("failed getting task")
		return err
	}

	if task.DeletedAt != nil {
		return h.Validate(c, http.StatusGone, echo.Map{"message": "task was deleted"})
	}

	if token.Subject() != task.CreatedBy.(*models.User).Id && !util.HasRoles(token, models.AdminRole.String(), models.SuperRole.String()) {
		return h.Validate(c, http.StatusForbidden, echo.Map{"message": "you don't have access"})
	}

	err = h.svc.Delete(ctx, token.Subject(), task)
	if err != nil {
		log.Error().Err(err).Msg("failed deleting task")
		return err
	}

	return h.Validate(c, http.StatusNoContent, nil)
}
