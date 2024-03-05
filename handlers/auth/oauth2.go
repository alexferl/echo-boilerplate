package auth

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
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
)

type GoogleUser struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func getOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  fmt.Sprintf("%s/oauth2/callback", viper.GetString(config.BaseURL)),
		ClientID:     viper.GetString(config.OAuth2ClientId),
		ClientSecret: viper.GetString(config.OAuth2ClientSecret),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func (h *Handler) OAuth2LogIn(c echo.Context) error {
	state, err := util.GenerateRandomString(80)
	if err != nil {
		return fmt.Errorf("oauth2: failed generating state: %v", err)
	}

	url := getOAuth2Config().AuthCodeURL(state)
	opts := &util.CookieOptions{
		Name:     "state",
		Value:    state,
		Path:     "/oauth2/callback",
		SameSite: http.SameSiteLaxMode, // needs to be Lax since it's across domains
		HttpOnly: true,
		MaxAge:   600,
	}
	c.SetCookie(util.NewCookie(opts))

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) OAuth2Callback(c echo.Context) error {
	response, err := callback(c)
	if err != nil {
		log.Error().Err(err).Msg("oauth2: failed callback")
		return err
	}

	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msg("oauth2: failed reading body")
		return err
	}

	if response.StatusCode != http.StatusOK {
		log.Error().Msgf("oauth2: response code was: %d body: %s", response.StatusCode, b)
		return c.JSON(http.StatusUnauthorized, echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "failed to log in",
		})
	}

	googleUser := &GoogleUser{}
	err = json.Unmarshal(b, googleUser)
	if err != nil {
		log.Error().Err(err).Msg("oauth2: failed unmarshalling body")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"email", googleUser.Email}}
	res, err := h.Mapper.FindOne(ctx, filter, &users.User{})
	if err != nil {
		if !errors.Is(err, data.ErrNoDocuments) {
			log.Error().Err(err).Msg("oauth2: failed getting user")
			return err
		}
	}
	var access, refresh []byte

	if res == nil {
		// TODO: username?
		newUser := users.NewUser(googleUser.Email, googleUser.Email)
		access, refresh, err = newUser.Login()
		if err != nil {
			log.Error().Err(err).Msg("oauth2: failed generating tokens")
			return err
		}

		newUser.Create(newUser.Id)

		_, err = h.Mapper.InsertOne(ctx, newUser, nil)
		if err != nil {
			log.Error().Err(err).Msg("oauth2: failed inserting user")
			return err
		}
	} else {
		user := res.(*users.User)
		access, refresh, err = user.Login()
		if err != nil {
			log.Error().Err(err).Msg("oauth2: failed generating tokens")
			return err
		}

		_, err = h.Mapper.UpdateOneById(ctx, user.Id, user, nil)
		if err != nil {
			log.Error().Err(err).Msg("oauth2: failed updating user")
			return err
		}
	}

	stateOpts := &util.CookieOptions{
		Name:     "state",
		Value:    "",
		Path:     "/oauth2/callback",
		SameSite: http.SameSiteLaxMode, // needs to be Lax since it's across domains
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(util.NewCookie(stateOpts))
	util.SetTokenCookies(c, access, refresh)

	resp := &TokenResponse{
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token, err := getOAuth2Config().Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed exchanging autorization code: %v", err)
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %v", err)
	}

	return resp, nil
}