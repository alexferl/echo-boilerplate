package users

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/util"
)

type AuthLogInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func (h *Handler) AuthLogIn(c echo.Context) error {
	body := &AuthLogInRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"email", body.Email}}
	res, err := h.Mapper.FindOne(ctx, filter, &User{})
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "invalid email or password"})
		}
		log.Error().Err(err).Msg("failed getting user")
		return err
	}

	user := res.(*User)
	err = user.ValidatePassword(body.Password)
	if err != nil {
		return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "invalid email or password"})
	}

	access, refresh, err := user.Login()
	if err != nil {
		log.Error().Err(err).Msg("failed generating tokens")
		return err
	}

	_, err = h.Mapper.UpdateOneById(ctx, user.Id, user, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
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
