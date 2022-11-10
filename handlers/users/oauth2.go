package users

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	libHttp "github.com/alexferl/golib/http/handler"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/alexferl/echo-boilerplate/config"
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
		RedirectURL:  fmt.Sprintf("%s/oauth2/callback", viper.GetString(config.BaseUrl)),
		ClientID:     viper.GetString(config.OAuth2ClientId),
		ClientSecret: viper.GetString(config.OAuth2ClientSecret),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func (h *Handler) OAuth2Login(c echo.Context) error {
	state, err := util.GenerateRandomString(80)
	if err != nil {
		return fmt.Errorf("oauth2: failed to generate state: %v", err)
	}

	url := getOAuth2Config().AuthCodeURL(state)
	opts := &util.CookieOptions{
		Name:     "state",
		Value:    state,
		Path:     "/oauth2/callback",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		MaxAge:   600,
	}
	c.SetCookie(util.NewCookie(opts))

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *Handler) OAuth2Callback(c echo.Context) error {
	response, err := callback(c)
	if err != nil {
		return fmt.Errorf("oauth2: failed callback: %v", err)
	}

	defer response.Body.Close()
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("oauth2: failed to read body: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		c.Logger().Errorf("oauth2: response code was: %d body: %s", response.StatusCode, b)
		return libHttp.JSONError(c, http.StatusUnauthorized, "failed to log in")
	}

	googleUser := &GoogleUser{}
	err = json.Unmarshal(b, googleUser)
	if err != nil {
		return fmt.Errorf("oauth2: failed unmarshalling body: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"email", googleUser.Email}}
	result, err := h.Mapper.FindOne(ctx, filter, &User{})
	if err != nil {
		if err != ErrUserNotFound {
			return fmt.Errorf("oauth2: failed to get user: %v", err)
		}
	}
	var access, refresh []byte

	if result == nil {
		newUser := NewUser(googleUser.Email, googleUser.Email)
		access, refresh, err = newUser.Login()
		if err != nil {
			return fmt.Errorf("oauth2: failed to generate tokens: %v", err)
		}

		newUser.Create(newUser.Id)

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = h.Mapper.Insert(ctx, newUser)
		if err != nil {
			return fmt.Errorf("oauth2: failed to insert user: %v", err)
		}
	} else {
		user := result.(*User)
		access, refresh, err = user.Login()
		if err != nil {
			return fmt.Errorf("oauth2: failed to generate tokens: %v", err)
		}

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err = h.Mapper.UpdateById(ctx, user.Id, user, nil)
		if err != nil {
			return fmt.Errorf("oauth2: failed to update user: %v", err)
		}
	}

	stateOpts := &util.CookieOptions{
		Name:     "state",
		Value:    "",
		Path:     "/oauth2/callback",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(util.NewCookie(stateOpts))
	util.SetTokenCookies(c, string(access), string(refresh))

	resp := &TokenResponse{
		AccessToken:  string(access),
		ExpiresIn:    int64(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()),
		RefreshToken: string(refresh),
		TokenType:    "Bearer",
	}

	return h.Validate(c, http.StatusOK, resp)
}

func callback(c echo.Context) (*http.Response, error) {
	state := c.FormValue("state")
	code := c.FormValue("code")

	stateCooke, err := c.Cookie("state")
	if err != nil {
		return nil, fmt.Errorf("oauth2: cookie was empty")
	}

	if state != stateCooke.Value {
		return nil, fmt.Errorf("oauth2: state mismatch")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	token, err := getOAuth2Config().Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("oauth2: failed getting user info: %v", err)
	}

	return resp, nil
}
