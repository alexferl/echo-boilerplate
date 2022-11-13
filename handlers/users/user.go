package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type UserResponse struct {
	Id        string     `json:"id" bson:"id"`
	Username  string     `json:"username" bson:"username"`
	Email     string     `json:"email" bson:"email"`
	Name      string     `json:"name" bson:"name"`
	Bio       string     `json:"bio" bson:"bio"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
}

func (h *Handler) GetUser(c echo.Context) error {
	token := c.Get("token").(jwt.Token)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := h.Mapper.FindOneById(ctx, token.Subject(), &UserResponse{})
	if err != nil {
		return fmt.Errorf("failed getting user: %v", err)
	}

	return h.Validate(c, http.StatusOK, result)
}

type UpdateUserRequest struct {
	Email string `json:"email" bson:"email"`
	Name  string `json:"name" bson:"name"`
	Bio   string `json:"bio" bson:"bio"`
}

func (h *Handler) UpdateUser(c echo.Context) error {
	body := &UpdateUserRequest{}
	if err := c.Bind(body); err != nil {
		return err
	}

	token := c.Get("token").(jwt.Token)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := h.Mapper.FindOneById(ctx, token.Subject(), &User{})
	if err != nil {
		return fmt.Errorf("failed getting user: %v", err)
	}

	user := result.(*User)
	if body.Email != "" {
		user.Email = body.Email
	}

	if body.Name != "" {
		user.Name = body.Name
	}

	if body.Bio != "" {
		user.Bio = body.Bio
	}

	user.Update(user.Id)

	update, err := h.Mapper.UpdateById(ctx, user.Id, user, &UserResponse{})
	if err != nil {
		return fmt.Errorf("failed updating user: %v", err)
	}

	return h.Validate(c, http.StatusOK, update)
}
