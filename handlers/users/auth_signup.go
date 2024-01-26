package users

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexferl/echo-boilerplate/data"
)

type AuthSignUpRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Password string `json:"password"`
}

func (h *Handler) AuthSignUp(c echo.Context) error {
	body := &AuthSignUpRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"$or", bson.A{
		bson.D{{"username", body.Username}},
		bson.D{{"email", body.Email}},
	}}}
	exist, err := h.Mapper.FindOne(ctx, filter, &User{})
	if err != nil {
		if !errors.Is(err, data.ErrNoDocuments) {
			log.Error().Err(err).Msg("failed finding user")
			return err
		}
	}

	if exist != nil {
		return h.Validate(c, http.StatusConflict, echo.Map{"message": "email or username already in-use"})
	}

	newUser := NewUser(body.Email, body.Username)
	newUser.Name = body.Name
	newUser.Bio = body.Bio
	err = newUser.SetPassword(body.Password)
	if err != nil {
		log.Error().Err(err).Msg("failed setting password")
		return err
	}

	newUser.Create(newUser.Id)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	user, err := h.Mapper.FindOneAndUpdate(ctx, filter, newUser, &User{}, opts)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return h.Validate(c, http.StatusConflict, echo.Map{"message": "email or username already in-use"})
		}
		log.Error().Err(err).Msg("failed inserting new user")
		return err
	}

	return h.Validate(c, http.StatusOK, user.(*User).Response())
}
