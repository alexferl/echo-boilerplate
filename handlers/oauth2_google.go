package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/services"
	"github.com/alexferl/echo-boilerplate/util/cookie"
	"github.com/alexferl/echo-boilerplate/util/rand"
)

type GoogleUser struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func getOAuth2GoogleConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     viper.GetString(config.OAuth2GoogleClientId),
		ClientSecret: viper.GetString(config.OAuth2GoogleClientSecret),
		Endpoint:     google.Endpoint,
		RedirectURL:  fmt.Sprintf("%s/oauth2/google/callback", viper.GetString(config.BaseURL)),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	}
}

func (h *AuthHandler) oauth2GoogleLogin(c echo.Context) error {
	state, err := rand.GenerateRandomString(80)
	if err != nil {
		return fmt.Errorf("google: failed generating state: %v", err)
	}

	url := getOAuth2GoogleConfig().AuthCodeURL(state)
	opts := &cookie.Options{
		Name:     "state",
		Value:    state,
		Path:     "/oauth2/google/callback",
		SameSite: http.SameSiteLaxMode, // needs to be Lax since it's across domains
		HttpOnly: true,
		MaxAge:   600,
	}
	c.SetCookie(cookie.New(opts))

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

type OAuth2GoogleCallbackResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func (h *AuthHandler) oauth2GoogleCallback(c echo.Context) error {
	response, err := callback(c)
	if err != nil {
		log.Error().Err(err).Msg("google: failed callback")
		return err
	}

	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("google: failed reading body")
		return err
	}

	if response.StatusCode != http.StatusOK {
		log.Error().Msgf("google: response code was: %d body: %s", response.StatusCode, b)
		return c.JSON(http.StatusUnauthorized, echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "failed to log in",
		})
	}

	googleUser := &GoogleUser{}
	err = json.Unmarshal(b, googleUser)
	if err != nil {
		log.Error().Err(err).Msg("google: failed unmarshalling body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	res, err := h.svc.FindOneByEmailOrUsername(ctx, googleUser.Email, "")
	if err != nil {
		var se *services.Error
		if errors.As(err, &se) {
			if se.Kind != services.NotExist {
				log.Error().Err(err).Msg("failed getting user")
			}
		}
	}

	var access, refresh []byte

	if res == nil {
		newUser := models.NewUser(googleUser.Email, "")
		access, refresh, err = newUser.Login()
		if err != nil {
			log.Error().Err(err).Msg("google: failed generating tokens")
			return err
		}

		_, err = h.svc.Create(ctx, newUser)
		if err != nil {
			log.Error().Err(err).Msg("google: failed inserting user")
			return err
		}
	} else {
		user := res
		access, refresh, err = user.Login()
		if err != nil {
			log.Error().Err(err).Msg("google: failed generating tokens")
			return err
		}

		_, err = h.svc.Update(ctx, user.Id, user)
		if err != nil {
			log.Error().Err(err).Msg("google: failed updating user")
			return err
		}
	}

	stateOpts := &cookie.Options{
		Name:     "state",
		Value:    "",
		Path:     "/oauth2/callback",
		SameSite: http.SameSiteLaxMode, // needs to be Lax since it's across domains
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(cookie.New(stateOpts))
	cookie.SetToken(c, access, refresh)

	resp := &OAuth2GoogleCallbackResponse{
		AccessToken:  string(access),
		ExpiresIn:    int64(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()),
		RefreshToken: string(refresh),
		TokenType:    "Bearer",
	}

	return c.JSON(http.StatusOK, resp)
}

func callback(c echo.Context) (*http.Response, error) {
	state := c.FormValue("state")
	code := c.FormValue("code")

	stateCooke, err := c.Cookie("state")
	if err != nil {
		return nil, fmt.Errorf("cookie was empty")
	}

	if state != stateCooke.Value {
		return nil, fmt.Errorf("state mismatch")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()
	token, err := getOAuth2GoogleConfig().Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed exchanging autorization code: %v", err)
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %v", err)
	}

	return resp, nil
}
