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

func (h *Handler) GetUser(c echo.Context) error {
	id := c.Param("id_or_username")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"$or", bson.A{
		bson.D{{"id", id}},
		bson.D{{"username", id}},
	}}}
	res, err := h.Mapper.FindOne(ctx, filter, &User{})
	if errors.Is(err, data.ErrNoDocuments) {
		return h.Validate(c, http.StatusNotFound, echo.Map{"message": "user not found"})
	} else if err != nil {
		log.Error().Err(err).Msg("failed getting user")
	}

	user := res.(*User)
	if user.DeletedAt != nil {
		return h.Validate(c, http.StatusGone, echo.Map{"message": "user deleted"})
	}

	return h.Validate(c, http.StatusOK, user.Public())
}
