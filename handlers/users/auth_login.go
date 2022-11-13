package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/util"
)

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) AuthLogin(c echo.Context) error {
	body := &LoginPayload{}
	if err := c.Bind(body); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"email", body.Email}}
	result, err := h.Mapper.FindOne(ctx, filter, &User{})
	if err != nil {
		if err == ErrUserNotFound {
			return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "invalid email or password"})
		}
		return fmt.Errorf("failed getting user: %v", err)
	}

	user := result.(*User)
	err = user.CheckPassword(body.Password)
	if err != nil {
		return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "invalid email or password"})
	}

	access, refresh, err := user.Login()
	if err != nil {
		return fmt.Errorf("failed generating tokens: %v", err)
	}

	_, err = h.Mapper.UpdateById(ctx, user.Id, user, nil)
	if err != nil {
		return fmt.Errorf("failed updating user: %v", err)
	}

	util.SetTokenCookies(c, access, refresh)

	resp := &TokenResponse{
		AccessToken:  string(access),
		ExpiresIn:    int64(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()),
		RefreshToken: string(refresh),
		TokenType:    "Bearer",
	}

	return h.Validate(c, http.StatusOK, resp)
}
