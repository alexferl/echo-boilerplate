package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/labstack/echo/v4"
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
		return fmt.Errorf("failed to bind: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"$or", bson.A{
		bson.D{{"username", body.Username}},
		bson.D{{"email", body.Email}},
	}}}
	exist, err := h.Mapper.FindOne(ctx, filter, &UserResponse{})
	if err != nil {
		if err != ErrNoDocuments {
			return fmt.Errorf("failed to get newUser: %v", err)
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
		return fmt.Errorf("failed to set password: %v", err)
	}

	newUser.Create(newUser.Id)

	opts := options.FindOneAndUpdate().SetUpsert(true)
	user, err := h.Mapper.Upsert(ctx, filter, newUser, &UserResponse{}, opts)
	if err != nil {
		return fmt.Errorf("failed to insert newUser: %v", err)
	}

	return h.Validate(c, http.StatusOK, user)
}
