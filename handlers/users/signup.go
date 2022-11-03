package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/labstack/echo/v4"
)

type SignUpPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (h *Handler) SignUp(c echo.Context) error {
	body := &SignUpPayload{}
	if err := c.Bind(body); err != nil {
		return fmt.Errorf("failed to bind: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"$or", bson.A{
		bson.D{{"username", body.Username}},
		bson.D{{"email", body.Email}},
	}}}
	exist, err := h.Mapper.FindOne(ctx, filter, &SignUpResponse{})
	if err != nil {
		if err != ErrUserNotFound {
			return fmt.Errorf("failed to get user: %v", err)
		}
	}

	if exist != nil {
		return h.Validate(c, http.StatusConflict, echo.Map{"message": "email already in-use"})
	}

	user := NewUser(body.Email, body.Username)
	err = user.SetPassword(body.Password)
	if err != nil {
		return fmt.Errorf("failed to set password: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user.Create(user.Id)

	err = h.Mapper.Insert(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	u := SignUpResponse{
		Id:       user.Id,
		Email:    user.Email,
		Username: user.Username,
	}

	return h.Validate(c, http.StatusOK, u)
}
