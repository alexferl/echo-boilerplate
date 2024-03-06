package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/util"
)

type AuthHandler struct {
	*openapi.Handler
	svc UserService
}

func NewAuthHandler(openapi *openapi.Handler, svc UserService) *AuthHandler {
	return &AuthHandler{
		Handler: openapi,
		svc:     svc,
	}
}

func (h *AuthHandler) Register(s *server.Server) {
	s.Add(http.MethodPost, "/auth/login", h.login)
	s.Add(http.MethodPost, "/auth/logout", h.logout)
	s.Add(http.MethodPost, "/auth/refresh", h.refresh)
	s.Add(http.MethodPost, "/auth/signup", h.signup)
	s.Add(http.MethodGet, "/auth/token", h.token)
	// s.Add(http.MethodGet, "/oauth2/callback", h.OAuth2Callback)
	// s.Add(http.MethodGet, "/oauth2/login", h.OAuth2LogIn)
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
	Username string `json:"username,omitempty"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func (h *AuthHandler) login(c echo.Context) error {
	body := &LoginRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.FindOneByEmailOrUsername(ctx, body.Email, body.Username)
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "invalid email or password"})
		}
		log.Error().Err(err).Msg("failed finding user")
		return err
	}

	err = user.ValidatePassword(body.Password)
	if err != nil {
		return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "invalid email or password"})
	}

	access, refresh, err := user.Login()
	if err != nil {
		log.Error().Err(err).Msg("failed generating tokens")
		return err
	}

	_, err = h.svc.Update(ctx, "", user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	if viper.GetBool(config.CookiesEnabled) {
		util.SetTokenCookies(c, access, refresh)
	}

	resp := &LoginResponse{
		AccessToken:  string(access),
		ExpiresIn:    int64(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()),
		RefreshToken: string(refresh),
		TokenType:    "Bearer",
	}

	return h.Validate(c, http.StatusOK, resp)
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) logout(c echo.Context) error {
	token := c.Get("refresh_token").(jwt.Token)
	encodedToken := c.Get("refresh_token_encoded").(string)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, token.Subject())
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "token not found"})
		}
		log.Error().Err(err).Msg("failed getting user")
		return err
	}

	if err = user.ValidateRefreshToken(encodedToken); err != nil {
		return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "token mismatch"})
	}

	user.Logout()

	_, err = h.svc.Update(ctx, "", user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	util.SetExpiredTokenCookies(c)

	return h.Validate(c, http.StatusNoContent, nil)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func (h *AuthHandler) refresh(c echo.Context) error {
	token := c.Get("refresh_token").(jwt.Token)
	encodedToken := c.Get("refresh_token_encoded").(string)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	user, err := h.svc.Read(ctx, token.Subject())
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "token not found"})
		}
		log.Error().Err(err).Msg("failed getting user")
		return err
	}

	if err = user.ValidateRefreshToken(encodedToken); err != nil {
		return h.Validate(c, http.StatusUnauthorized, echo.Map{"message": "token mismatch"})
	}

	access, refresh, err := user.Refresh()
	if err != nil {
		log.Error().Err(err).Msg("failed generating tokens")
		return err
	}

	_, err = h.svc.Update(ctx, "", user)
	if err != nil {
		log.Error().Err(err).Msg("failed updating user")
		return err
	}

	if viper.GetBool(config.CookiesEnabled) {
		util.SetTokenCookies(c, access, refresh)
	}

	resp := &RefreshResponse{
		AccessToken:  string(access),
		ExpiresIn:    int64(viper.GetDuration(config.JWTAccessTokenExpiry).Seconds()),
		RefreshToken: string(refresh),
		TokenType:    "Bearer",
	}

	return h.Validate(c, http.StatusOK, resp)
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
	Password string `json:"password"`
}

func (h *AuthHandler) signup(c echo.Context) error {
	body := &SignUpRequest{}
	if err := c.Bind(body); err != nil {
		log.Error().Err(err).Msg("failed binding body")
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	res, err := h.svc.FindOneByEmailOrUsername(ctx, body.Email, body.Username)
	if err != nil {
		if !errors.Is(err, data.ErrNoDocuments) {
			log.Error().Err(err).Msg("failed getting user")
			return err
		}
	}

	if res != nil {
		return h.Validate(c, http.StatusConflict, echo.Map{"message": "email or username already in-use"})
	}

	user := models.NewUser(body.Email, body.Username)
	user.Name = body.Name
	user.Bio = body.Bio
	err = user.SetPassword(body.Password)
	if err != nil {
		log.Error().Err(err).Msg("failed setting password")
		return err
	}

	user.Create(user.Id)

	res, err = h.svc.Create(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return h.Validate(c, http.StatusConflict, echo.Map{"message": "email or username already in-use"})
		}
		log.Error().Err(err).Msg("failed inserting new user")
		return err
	}

	return h.Validate(c, http.StatusOK, res.Response())
}

type TokenResponse struct {
	Exp   time.Time `json:"exp"`
	Iat   time.Time `json:"iat"`
	Iss   string    `json:"iss"`
	Nbf   time.Time `json:"nbf"`
	Roles []string  `json:"roles"`
	Sub   string    `json:"sub"`
	Type  string    `json:"type"`
}

func (h *AuthHandler) token(c echo.Context) error {
	token := c.Get("token").(jwt.Token)

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	m, err := token.AsMap(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed converting token to map")
		return err
	}

	return h.Validate(c, http.StatusOK, m)
}
