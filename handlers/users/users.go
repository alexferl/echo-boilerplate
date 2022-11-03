package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ShortUser struct {
	Id       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
}

type UsersResponse struct {
	Users []*ShortUser `json:"users"`
}

func (h *Handler) Users(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := h.Mapper.Find(ctx, nil, []*ShortUser{})
	if err != nil {
		return fmt.Errorf("failed getting users: %v", err)
	}

	resp := &UsersResponse{Users: result.([]*ShortUser)}

	return h.Validate(c, http.StatusOK, resp)
}
