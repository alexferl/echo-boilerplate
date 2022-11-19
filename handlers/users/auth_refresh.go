package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/spf13/viper"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/util"
)

type AuthRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) AuthRefresh(c echo.Context) error {
	token := c.Get("refresh_token").(jwt.Token)
	encodedToken := c.Get("refresh_token_encoded").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := h.Mapper.FindOneById(ctx, token.Subject(), &User{})
	if err != nil {
		return fmt.Errorf("failed getting user: %v", err)
	}

	user := result.(*User)
	if err = user.ValidateRefreshToken(encodedToken); err != nil {
		return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "Token mismatch"})
	}

	access, refresh, err := user.Refresh()
	if err != nil {
		return fmt.Errorf("failed generating tokens: %v", err)
	}

	_, err = h.Mapper.UpdateById(ctx, token.Subject(), user, nil)
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
