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

type UserService interface {
	Create(ctx context.Context, model *models.User) (*models.User, error)
	Read(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id string, model *models.User) (*models.User, error)
	Delete(ctx context.Context, id string, model *models.User) error
	Find(ctx context.Context, params *models.UserSearchParams) (int64, models.Users, error)
	FindOneByEmailOrUsername(ctx context.Context, email string, username string) (*models.User, error)
}

type UserHandler struct {
	*openapi.Handler
	svc UserService
}

func NewUserHandler(openapi *openapi.Handler, svc UserService) *UserHandler {
	return &UserHandler{
		Handler: openapi,
		svc:     svc,
	}
}

func (h *UserHandler) Register(s *server.Server) {
	s.Add(http.MethodGet, "/me", h.getCurrentUser)
	s.Add(http.MethodPatch, "/me", h.updateCurrentUser)
	s.Add(http.MethodGet, "/users/:username", h.get)
	s.Add(http.MethodPatch, "/users/:username", h.update)
	s.Add(http.MethodPut, "/users/:username/ban", h.ban)
	s.Add(http.MethodDelete, "/users/:username/ban", h.unban)
	s.Add(http.MethodPut, "/users/:username/lock", h.lock)
	s.Add(http.MethodDelete, "/users/:username/lock", h.unlock)
	s.Add(http.MethodPut, "/users/:username/roles/:role", h.addRole)
	s.Add(http.MethodDelete, "/users/:username/roles/:role", h.removeRole)
	s.Add(http.MethodGet, "/users", h.list)
}

func (h *UserHandler) getCurrentUser(c echo.Context) error {
	currentUser := c.Get("user").(*models.User)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, currentUser.Id)
	if err != nil {
		log.Error().Err(err).Msg("failed getting user")
		return err
	}

	return h.Validate(c, http.StatusOK, user.Response())
}

type UpdateCurrentUserRequest struct {
	Name *string `json:"name,omitempty"`
	Bio  *string `json:"bio,omitempty"`
}

func (h *UserHandler) updateCurrentUser(c echo.Context) error {
	currentUser := c.Get("user").(*models.User)

	body := &UpdateCurrentUserRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, currentUser.Id)
	if err != nil {
		log.Error().Err(err).Msg("failed getting user")
		return err
	}

	if body.Name != nil {
		user.Name = *body.Name
	}

	if body.Bio != nil {
		user.Bio = *body.Bio
	}

	res, err := h.svc.Update(ctx, currentUser.Id, user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusOK, res.Response())
}

func (h *UserHandler) get(c echo.Context) error {
	id := c.Param("username")
	currentUser := c.Get("user").(*models.User)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, id)
	if err != nil {
		sErr := h.readUser(c, err)
		if sErr != nil {
			return sErr()
		}
	}

	if currentUser.HasRoleOrHigher(models.AdminRole) {
		return h.Validate(c, http.StatusOK, user.AdminResponse())
	}

	return h.Validate(c, http.StatusOK, user.Response())
}

type UpdateUserRequest struct {
	Name *string `json:"name,omitempty"`
	Bio  *string `json:"bio,omitempty"`
}

func (h *UserHandler) update(c echo.Context) error {
	id := c.Param("username")
	currentUser := c.Get("user").(*models.User)

	body := &UpdateUserRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, id)
	if err != nil {
		sErr := h.readUser(c, err)
		if sErr != nil {
			return sErr()
		}
	}

	if body.Name != nil {
		user.Name = *body.Name
	}

	if body.Bio != nil {
		user.Bio = *body.Bio
	}

	res, err := h.svc.Update(ctx, currentUser.Id, user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusOK, res.AdminResponse())
}

func (h *UserHandler) ban(c echo.Context) error {
	id := c.Param("username")
	currentUser := c.Get("user").(*models.User)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, id)
	if err != nil {
		sErr := h.readUser(c, err)
		if sErr != nil {
			return sErr()
		}
	}

	err = user.Ban(currentUser)
	if err != nil {
		mErr := h.checkModelErr(c, err, "banning")
		if mErr != nil {
			return mErr()
		}
	}

	_, err = h.svc.Update(ctx, currentUser.Id, user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusNoContent, nil)
}

func (h *UserHandler) unban(c echo.Context) error {
	id := c.Param("username")
	currentUser := c.Get("user").(*models.User)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, id)
	if err != nil {
		sErr := h.readUser(c, err)
		if sErr != nil {
			return sErr()
		}
	}

	err = user.Unban(currentUser)
	if err != nil {
		mErr := h.checkModelErr(c, err, "unbanning")
		if mErr != nil {
			return mErr()
		}
	}

	_, err = h.svc.Update(ctx, currentUser.Id, user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusNoContent, nil)
}

func (h *UserHandler) lock(c echo.Context) error {
	id := c.Param("username")
	currentUser := c.Get("user").(*models.User)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, id)
	if err != nil {
		sErr := h.readUser(c, err)
		if sErr != nil {
			return sErr()
		}
	}

	err = user.Lock(currentUser)
	if err != nil {
		mErr := h.checkModelErr(c, err, "locking")
		if mErr != nil {
			return mErr()
		}
	}

	_, err = h.svc.Update(ctx, currentUser.Id, user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusNoContent, nil)
}

func (h *UserHandler) unlock(c echo.Context) error {
	id := c.Param("username")
	currentUser := c.Get("user").(*models.User)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, id)
	if err != nil {
		sErr := h.readUser(c, err)
		if sErr != nil {
			return sErr()
		}
	}

	err = user.Unlock(currentUser)
	if err != nil {
		mErr := h.checkModelErr(c, err, "locking")
		if mErr != nil {
			return mErr()
		}
	}

	_, err = h.svc.Update(ctx, currentUser.Id, user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusNoContent, nil)
}

func (h *UserHandler) addRole(c echo.Context) error {
	id := c.Param("username")
	role := c.Param("role")
	currentUser := c.Get("user").(*models.User)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, id)
	if err != nil {
		sErr := h.readUser(c, err)
		if sErr != nil {
			return sErr()
		}
	}

	err = user.AddRole(currentUser, models.RolesMap[role])
	if err != nil {
		mErr := h.checkModelErr(c, err, "locking")
		if mErr != nil {
			return mErr()
		}
	}

	_, err = h.svc.Update(ctx, currentUser.Id, user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusNoContent, nil)
}

func (h *UserHandler) removeRole(c echo.Context) error {
	id := c.Param("username")
	role := c.Param("role")
	currentUser := c.Get("user").(*models.User)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, id)
	if err != nil {
		sErr := h.readUser(c, err)
		if sErr != nil {
			return sErr()
		}
	}

	err = user.RemoveRole(currentUser, models.RolesMap[role])
	if err != nil {
		mErr := h.checkModelErr(c, err, "locking")
		if mErr != nil {
			return mErr()
		}
	}

	_, err = h.svc.Update(ctx, currentUser.Id, user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusNoContent, nil)
}

func (h *UserHandler) list(c echo.Context) error {
	page, perPage, limit, skip := pagination.ParseParams(c)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	params := &models.UserSearchParams{
		Limit: limit,
		Skip:  skip,
	}
	count, users, err := h.svc.Find(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed getting tasks")
		return err
	}

	pagination.SetHeaders(c.Request(), c.Response().Header(), int(count), page, perPage)

	return h.Validate(c, http.StatusOK, users.AdminResponse())
}

func (h *UserHandler) readUser(c echo.Context, err error) func() error {
	var se *services.Error
	if errors.As(err, &se) {
		msg := echo.Map{"message": se.Message}
		if se.Kind == services.NotExist {
			return func() error { return h.Validate(c, http.StatusNotFound, msg) }
		} else if se.Kind == services.Deleted {
			return func() error { return h.Validate(c, http.StatusGone, msg) }
		}
	}
	log.Error().Err(err).Msg("failed getting user")
	return func() error { return err }
}

func (h *UserHandler) checkModelErr(c echo.Context, err error, action string) func() error {
	var me *models.Error
	if errors.As(err, &me) {
		msg := echo.Map{"message": me.Message}
		if me.Kind == models.Conflict {
			return func() error { return h.Validate(c, http.StatusConflict, msg) }
		} else if me.Kind == models.Permission {
			return func() error { return h.Validate(c, http.StatusForbidden, msg) }
		}
	}
	log.Error().Err(err).Msgf("failed %s user", action)
	return func() error { return err }
}
