package users

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog"
)

func (h *Handler) GetCurrentUser(c echo.Context) error {
	token := c.Get("token").(jwt.Token)
	logger := c.Get("logger").(zerolog.Logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := h.Mapper.FindOneById(ctx, token.Subject(), &User{})
	if err != nil {
		logger.Error().Err(err).Msg("failed getting user")
		return err
	}

	return h.Validate(c, http.StatusOK, res.(*User).Response())
}

type UpdateCurrentUserRequest struct {
	Name string `json:"name" bson:"name"`
	Bio  string `json:"bio" bson:"bio"`
}

func (h *Handler) UpdateCurrentUser(c echo.Context) error {
	token := c.Get("token").(jwt.Token)
	logger := c.Get("logger").(zerolog.Logger)

	body := &UpdateCurrentUserRequest{}
	if err := c.Bind(body); err != nil {
		logger.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := h.Mapper.FindOneById(ctx, token.Subject(), &User{})
	if err != nil {
		logger.Error().Err(err).Msg("failed getting user")
		return err
	}

	user := res.(*User)
	if body.Name != "" {
		user.Name = body.Name
	}

	if body.Bio != "" {
		user.Bio = body.Bio
	}

	user.Update(user.Id)

	update, err := h.Mapper.FindOneByIdAndUpdate(ctx, user.Id, user, &User{})
	if err != nil {
		logger.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusOK, update.(*User).Response())
}
