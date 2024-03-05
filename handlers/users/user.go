package users

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/util"
)

func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id_or_username")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, errResp := h.getUser(ctx, c, id)
	if errResp != nil {
		return errResp()
	}

	return h.Validate(c, http.StatusOK, user.Public())
}

type UpdateUserRequest struct {
	Name *string `json:"name" bson:"name"`
	Bio  *string `json:"bio" bson:"bio"`
}

func (h *Handler) UpdateUser(c echo.Context) error {
	id := c.Param("id_or_username")
	logger := c.Get("logger").(zerolog.Logger)
	token := c.Get("token").(jwt.Token)

	body := &UpdateUserRequest{}
	if err := c.Bind(body); err != nil {
		logger.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user, errResp := h.getUser(ctx, c, id)
	if errResp != nil {
		return errResp()
	}

	if body.Name != nil {
		user.Name = *body.Name
	}

	if body.Bio != nil {
		user.Bio = *body.Bio
	}

	user.Update(token.Subject())

	update, err := h.Mapper.FindOneByIdAndUpdate(ctx, user.Id, user, &User{})
	if err != nil {
		logger.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusOK, update.(*User).AdminResponse())
}

type UpdateUserStatusRequest struct {
	IsBanned *bool `json:"is_banned"`
	IsLocked *bool `json:"is_locked"`
}

func (h *Handler) UpdateUserStatus(c echo.Context) error {
	id := c.Param("id_or_username")
	logger := c.Get("logger").(zerolog.Logger)
	token := c.Get("token").(jwt.Token)

	body := &UpdateUserStatusRequest{}
	if err := c.Bind(body); err != nil {
		logger.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user, errResp := h.getUser(ctx, c, id)
	if errResp != nil {
		return errResp()
	}

	if user.Id == token.Subject() {
		return c.JSON(http.StatusForbidden, echo.Map{"message": "you cannot update your own status"})
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

	update, err := h.Mapper.FindOneByIdAndUpdate(ctx, user.Id, user, &User{})
	if err != nil {
		logger.Error().Err(err).Msg("failed updating user")
		return err
	}

	return h.Validate(c, http.StatusOK, update.(*User).AdminResponse())
}

func (h *Handler) getUser(ctx context.Context, c echo.Context, id string) (*User, func() error) {
	logger, ok := c.Get("logger").(zerolog.Logger)
	if !ok {
		logger = log.Logger
	}

	filter := bson.D{{"$or", bson.A{
		bson.D{{"id", id}},
		bson.D{{"username", id}},
	}}}
	res, err := h.Mapper.FindOne(ctx, filter, &User{})
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return nil, util.WrapErr(h.Validate(c, http.StatusNotFound, echo.Map{"message": "user not found"}))
		}
		logger.Error().Err(err).Msg("failed getting user")
		return nil, util.WrapErr(err)
	}

	user := res.(*User)
	if user.DeletedAt != nil {
		return nil, util.WrapErr(h.Validate(c, http.StatusGone, echo.Map{"message": "user was deleted"}))
	}

	return user, nil
}
