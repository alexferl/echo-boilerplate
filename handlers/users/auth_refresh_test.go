package users_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
)

func TestHandler_AuthRefresh_200_Cookie(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	_, refresh, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewRefreshTokenCookie(refresh))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			user,
			nil,
		).
		On(
			"UpdateById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			nil,
			nil,
		)

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	expected := 2
	if viper.GetBool(config.CSRFEnabled) {
		expected = 3
	}
	if assert.Equal(t, expected, len(resp.Result().Cookies())) {
		cookies := 0
		for _, c := range resp.Result().Cookies() {
			if c.Name == viper.GetString(config.JWTAccessTokenCookieName) {
				cookies++
			}
			if c.Name == viper.GetString(config.JWTRefreshTokenCookieName) {
				cookies++
			}
			if c.Name == viper.GetString(config.CSRFCookieName) {
				cookies++
			}
		}
		assert.Equal(t, expected, cookies)
	}
	assert.Contains(t, resp.Body.String(), "access_token")
	assert.Contains(t, resp.Body.String(), "expires_in")
	assert.Contains(t, resp.Body.String(), "refresh_token")
	assert.Contains(t, resp.Body.String(), "token_type")
}

func TestHandler_AuthRefresh_400_Cookie_Missing(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Request malformed")
}

func TestHandler_AuthRefresh_401_Cookie_Invalid(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewRefreshTokenCookie([]byte("invalid")))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Token invalid")
}

func TestHandler_AuthRefresh_401_Cookie_Mismatch(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	_, refresh, err := user.Login()
	assert.NoError(t, err)
	user.Logout()

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewRefreshTokenCookie(refresh))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			user,
			nil,
		)

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Token mismatch")
}

func TestHandler_AuthRefresh_200_Token(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	_, refresh, err := user.Login()
	assert.NoError(t, err)

	payload := &users.AuthRefreshRequest{
		RefreshToken: string(refresh),
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			user,
			nil,
		).
		On(
			"UpdateById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			nil,
			nil,
		)

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	expected := 2
	if viper.GetBool(config.CSRFEnabled) {
		expected = 3
	}
	if assert.Equal(t, expected, len(resp.Result().Cookies())) {
		cookies := 0
		for _, c := range resp.Result().Cookies() {
			if c.Name == viper.GetString(config.JWTAccessTokenCookieName) {
				cookies++
			}
			if c.Name == viper.GetString(config.JWTRefreshTokenCookieName) {
				cookies++
			}
			if c.Name == viper.GetString(config.CSRFCookieName) {
				cookies++
			}
		}
		assert.Equal(t, expected, cookies)
	}
	assert.Contains(t, resp.Body.String(), "access_token")
	assert.Contains(t, resp.Body.String(), "expires_in")
	assert.Contains(t, resp.Body.String(), "refresh_token")
	assert.Contains(t, resp.Body.String(), "token_type")
}

func TestHandler_AuthRefresh_400_Token_Missing(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Request malformed")
}

func TestHandler_AuthRefresh_401_Token_Invalid(t *testing.T) {
	_, s := getMapperAndServer(t)

	payload := &users.AuthRefreshRequest{
		RefreshToken: "invalid",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Token invalid")
}

func TestHandler_AuthRefresh_401_Token_Mismatch(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	_, refresh, err := user.Login()
	assert.NoError(t, err)
	user.Logout()

	payload := &users.AuthRefreshRequest{
		RefreshToken: string(refresh),
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			user,
			nil,
		)

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Token mismatch")
}
