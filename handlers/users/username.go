package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type GetUsernameResponse struct {
	Id        string     `json:"id"`
	Username  string     `json:"username"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	DeletedAt *time.Time `json:"-" bson:"deleted_at"`
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
	result, err := h.Mapper.FindOne(ctx, filter, &GetUsernameResponse{})
	if err == ErrUserNotFound {
		return h.Validate(c, http.StatusNotFound, echo.Map{"message": "user not found"})
	} else if err != nil {
		return fmt.Errorf("failed getting username: %v", err)
	}

	user := result.(*GetUsernameResponse)
	if user.DeletedAt != nil {
		return h.Validate(c, http.StatusGone, echo.Map{"message": "user deleted"})
	}

	return h.Validate(c, http.StatusOK, user)
}
