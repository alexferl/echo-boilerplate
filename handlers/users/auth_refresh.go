package users

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/util"
)

type AuthRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) AuthRefresh(c echo.Context) error {
	token := c.Get("refresh_token").(jwt.Token)
	encodedToken := c.Get("refresh_token_encoded").(string)
	logger := c.Get("logger").(zerolog.Logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := h.Mapper.FindOneById(ctx, token.Subject(), &User{})
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "Token not found"})
		}
		logger.Error().Err(err).Msg("failed getting user")
		return err
	}

	user := res.(*User)
	if err = user.ValidateRefreshToken(encodedToken); err != nil {
		return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "Token mismatch"})
	}

	access, refresh, err := user.Refresh()
	if err != nil {
		logger.Error().Err(err).Msg("failed generating tokens")
		return err
	}

	_, err = h.Mapper.UpdateOneById(ctx, token.Subject(), user, nil)
	if err != nil {
		logger.Error().Err(err).Msg("failed updating user")
		return err
	}

	if viper.GetBool(config.CookiesEnabled) {
		util.SetTokenCookies(c, access, refresh)
	}

	resp := &TokenResponse{
		AccessToken:  string(access),
		ExpiresIn:    int64(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()),
		RefreshToken: string(refresh),
		TokenType:    "Bearer",
	}

	return h.Validate(c, http.StatusOK, resp)
}
