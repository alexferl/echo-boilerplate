package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/alexferl/echo-boilerplate/util"
	"github.com/labstack/echo/v4"
)

type RevokePayload struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) AuthRevoke(c echo.Context) error {
	body := &RevokePayload{}
	if err := c.Bind(body); err != nil {
		return err
	}

	refreshToken := body.RefreshToken
	if refreshToken == "" {
		cookie, err := c.Cookie("refresh_token")
		if err != nil {
			return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "token missing"})
		}
		refreshToken = cookie.Value
	}

	token, err := util.ParseToken(refreshToken)
	if err != nil {
		return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "token invalid"})
	}

	hashedToken := util.HashToken([]byte(refreshToken))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := h.Mapper.FindOneById(ctx, token.Subject(), &User{})
	if err != nil {
		return fmt.Errorf("failed getting user: %v", err)
	}

	user := result.(*User)

	if user.RefreshTokenHash != hashedToken {
		return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "token mismatch"})
	}

	user.Logout()
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = h.Mapper.UpdateById(ctx, token.Subject(), user, nil)
	if err != nil {
		return fmt.Errorf("failed updating user: %v", err)
	}

	util.SetExpiredTokenCookies(c)

	return h.Validate(c, http.StatusNoContent, nil)
}
