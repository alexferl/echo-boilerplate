package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	jwx "github.com/lestrrat-go/jwx/v2/jwt"
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
	s.Add(http.MethodPut, "/me", h.updateCurrentUser)
	s.Add(http.MethodGet, "/users/:id", h.get)
	s.Add(http.MethodPut, "/users/:id", h.update)
	s.Add(http.MethodPut, "/users/:id/status", h.updateStatus)
	s.Add(http.MethodGet, "/users", h.list)
}

func (h *UserHandler) getCurrentUser(c echo.Context) error {
	token := c.Get("token").(jwx.Token)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, token.Subject())
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
	token := c.Get("token").(jwx.Token)

	body := &UpdateCurrentUserRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, token.Subject())
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

	res, err := h.svc.Update(ctx, token.Subject(), user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusOK, res.Response())
}

func (h *UserHandler) get(c echo.Context) error {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, id)
	if err != nil {
		var se *services.Error
		if errors.As(err, &se) {
			if se.Kind == services.NotExist {
				return h.Validate(c, http.StatusNotFound, echo.Map{"message": se.Message})
			} else if se.Kind == services.Deleted {
				return h.Validate(c, http.StatusGone, echo.Map{"message": se.Message})
			}
		}
		log.Error().Err(err).Msg("failed getting user")
		return err
	}

	return h.Validate(c, http.StatusOK, user.Response())
}

type UpdateUserRequest struct {
	Name *string `json:"name,omitempty"`
	Bio  *string `json:"bio,omitempty"`
}

func (h *UserHandler) update(c echo.Context) error {
	id := c.Param("id")
	token := c.Get("token").(jwx.Token)

	body := &UpdateUserRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, id)
	if err != nil {
		var se *services.Error
		if errors.As(err, &se) {
			if se.Kind == services.NotExist {
				return h.Validate(c, http.StatusNotFound, echo.Map{"message": se.Message})
			} else if se.Kind == services.Deleted {
				return h.Validate(c, http.StatusGone, echo.Map{"message": se.Message})
			}
		}
		log.Error().Err(err).Msg("failed getting user")
		return err
	}

	if body.Name != nil {
		user.Name = *body.Name
	}

	if body.Bio != nil {
		user.Bio = *body.Bio
	}

	res, err := h.svc.Update(ctx, token.Subject(), user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusOK, res.Response())
}

type UpdateUserStatusRequest struct {
	IsBanned *bool `json:"is_banned,omitempty"`
	IsLocked *bool `json:"is_locked,omitempty"`
}

func (h *UserHandler) updateStatus(c echo.Context) error {
	id := c.Param("id")
	token := c.Get("token").(jwx.Token)

	body := &UpdateUserStatusRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, id)
	if err != nil {
		var se *services.Error
		if errors.As(err, &se) {
			if se.Kind == services.NotExist {
				return h.Validate(c, http.StatusNotFound, echo.Map{"message": se.Message})
			} else if se.Kind == services.Deleted {
				return h.Validate(c, http.StatusGone, echo.Map{"message": se.Message})
			}
		}
		log.Error().Err(err).Msg("failed getting user")
		return err
	}

	if user.Id == token.Subject() {
		return c.JSON(http.StatusConflict, echo.Map{"message": "you cannot update your own status"})
	}

	if body.IsBanned != nil {
		banned := *body.IsBanned
		if banned {
			user.Ban(token.Subject())
		} else {
			user.Unban(token.Subject())
		}
	}

	if body.IsLocked != nil {
		locked := *body.IsLocked
		if locked {
			user.Lock(token.Subject())
		} else {
			user.Unlock(token.Subject())
		}
	}

	res, err := h.svc.Update(ctx, token.Subject(), user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusOK, res.Response())
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

	return h.Validate(c, http.StatusOK, users.Public())
}
