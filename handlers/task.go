package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/services"
	"github.com/alexferl/echo-boilerplate/util/pagination"
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
	s.Add(http.MethodPatch, "/tasks/:id", h.update)
	s.Add(http.MethodPut, "/tasks/:id/transition", h.transition)
	s.Add(http.MethodDelete, "/tasks/:id", h.delete)
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}

func (h *TaskHandler) create(c echo.Context) error {
	currentUser := c.Get("user").(*models.User)

	body := &CreateTaskRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	model := models.NewTask()
	model.Title = body.Title

	task, err := h.svc.Create(ctx, currentUser.Id, model)
	if err != nil {
		log.Error().Err(err).Msg("failed creating task")
		return err
	}

	return h.Validate(c, http.StatusOK, task.Response())
}

func (h *TaskHandler) list(c echo.Context) error {
	page, perPage, limit, skip := pagination.ParseParams(c)

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

	pagination.SetHeaders(c.Request(), c.Response().Header(), int(count), page, perPage)

	return h.Validate(c, http.StatusOK, tasks.Response())
}

func (h *TaskHandler) get(c echo.Context) error {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	task, err := h.svc.Read(ctx, id)
	if err != nil {
		sErr := h.readTask(c, err)
		if sErr != nil {
			return sErr()
		}
	}

	return h.Validate(c, http.StatusOK, task.Response())
}

type UpdateTaskRequest struct {
	Title *string `json:"title"`
}

func (h *TaskHandler) update(c echo.Context) error {
	id := c.Param("id")
	currentUser := c.Get("user").(*models.User)

	body := &UpdateTaskRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	task, err := h.svc.Read(ctx, id)
	if err != nil {
		sErr := h.readTask(c, err)
		if sErr != nil {
			return sErr()
		}
	}

	if currentUser.Id != task.CreatedBy.(*models.User).Id && !currentUser.HasRoleOrHigher(models.AdminRole) {
		return h.Validate(c, http.StatusForbidden, echo.Map{"message": "you don't have access"})
	}

	if body.Title != nil {
		task.Title = *body.Title
	}

	res, err := h.svc.Update(ctx, currentUser.Id, task)
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
	currentUser := c.Get("user").(*models.User)

	body := &TransitionTaskRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	task, err := h.svc.Read(ctx, id)
	if err != nil {
		sErr := h.readTask(c, err)
		if sErr != nil {
			return sErr()
		}
	}

	if *body.Completed != task.Completed {
		if *body.Completed {
			task.Complete(currentUser.Id)
		} else {
			task.Incomplete()
		}
	}

	res, err := h.svc.Update(ctx, currentUser.Id, task)
	if err != nil {
		log.Error().Err(err).Msg("failed updating task")
		return err
	}

	return h.Validate(c, http.StatusOK, res.Response())
}

func (h *TaskHandler) delete(c echo.Context) error {
	id := c.Param("id")
	currentUser := c.Get("user").(*models.User)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	task, err := h.svc.Read(ctx, id)
	if err != nil {
		sErr := h.readTask(c, err)
		if sErr != nil {
			return sErr()
		}
	}

	if currentUser.Id != task.CreatedBy.(*models.User).Id && !currentUser.HasRoleOrHigher(models.AdminRole) {
		return h.Validate(c, http.StatusForbidden, echo.Map{"message": "you don't have access"})
	}

	err = h.svc.Delete(ctx, currentUser.Id, task)
	if err != nil {
		log.Error().Err(err).Msg("failed deleting task")
		return err
	}

	return h.Validate(c, http.StatusNoContent, nil)
}

func (h *TaskHandler) readTask(c echo.Context, err error) func() error {
	var se *services.Error
	if errors.As(err, &se) {
		msg := echo.Map{"message": se.Message}
		if se.Kind == services.NotExist {
			return func() error { return h.Validate(c, http.StatusNotFound, msg) }
		} else if se.Kind == services.Deleted {
			return func() error { return h.Validate(c, http.StatusGone, msg) }
		}
	}
	log.Error().Err(err).Msg("failed getting task")
	return func() error { return err }
}
