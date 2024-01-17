package users

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexferl/echo-boilerplate/data"
)

type GetUsernameResponse struct {
	Id        string     `json:"id"`
	Href      string     `json:"href"`
	Username  string     `json:"username"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
}

func (h *Handler) GetUsername(c echo.Context) error {
	username := c.Param("username")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"$or", bson.A{
		bson.D{{"id", username}},
		bson.D{{"username", username}},
	}}}
	result, err := h.Mapper.FindOne(ctx, filter, &User{})
	if errors.Is(err, data.ErrNoDocuments) {
		return h.Validate(c, http.StatusNotFound, echo.Map{"message": "user not found"})
	} else if err != nil {
		log.Error().Err(err).Msg("failed getting user")
	}

	user := result.(*User)
	if user.DeletedAt != nil {
		return h.Validate(c, http.StatusGone, echo.Map{"message": "user deleted"})
	}

	resp := user.Response()

	return h.Validate(c, http.StatusOK, &GetUsernameResponse{
		Id:        resp.Id,
		Href:      resp.Href,
		Username:  resp.Username,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	})
}
