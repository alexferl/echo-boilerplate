package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Refresh_200_Cookie(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, refresh, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewAccessTokenCookie(string(access)))
	req.AddCookie(util.NewRefreshTokenCookie(string(refresh)))
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
	if assert.Equal(t, 2, len(resp.Result().Cookies())) {
		cookies := 0
		for _, c := range resp.Result().Cookies() {
			if c.Name == "access_token" {
				cookies++
			}
			if c.Name == "refresh_token" {
				cookies++
			}
		}
		assert.Equal(t, 2, cookies)
	}
	assert.Contains(t, resp.Body.String(), "access_token")
	assert.Contains(t, resp.Body.String(), "expires_in")
	assert.Contains(t, resp.Body.String(), "refresh_token")
	assert.Contains(t, resp.Body.String(), "token_type")
}

func TestHandler_Refresh_401_Cookie_Missing(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "token missing")
}

func TestHandler_Refresh_401_Cookie_Invalid(t *testing.T) {
	_, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewAccessTokenCookie(string(access)))
	req.AddCookie(util.NewRefreshTokenCookie("invalid"))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "token invalid")
}

func TestHandler_Refresh_401_Cookie_Mismatch(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, refresh, err := user.Login()
	assert.NoError(t, err)
	user.Logout()

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewAccessTokenCookie(string(access)))
	req.AddCookie(util.NewRefreshTokenCookie(string(refresh)))
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
	assert.Contains(t, resp.Body.String(), "token mismatch")
}

func TestHandler_Refresh_200_Token(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, refresh, err := user.Login()
	assert.NoError(t, err)

	payload := &users.RefreshPayload{
		RefreshToken: string(refresh),
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
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
	if assert.Equal(t, 2, len(resp.Result().Cookies())) {
		cookies := 0
		for _, c := range resp.Result().Cookies() {
			if c.Name == "access_token" {
				cookies++
			}
			if c.Name == "refresh_token" {
				cookies++
			}
		}
		assert.Equal(t, 2, cookies)
	}
	assert.Contains(t, resp.Body.String(), "access_token")
	assert.Contains(t, resp.Body.String(), "expires_in")
	assert.Contains(t, resp.Body.String(), "refresh_token")
	assert.Contains(t, resp.Body.String(), "token_type")
}

func TestHandler_Refresh_401_Token_Missing(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "token missing")
}

func TestHandler_Refresh_401_Token_Invalid(t *testing.T) {
	_, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &users.RefreshPayload{
		RefreshToken: "invalid",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "token invalid")
}

func TestHandler_Refresh_401_Token_Mismatch(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, refresh, err := user.Login()
	assert.NoError(t, err)
	user.Logout()

	payload := &users.RefreshPayload{
		RefreshToken: string(refresh),
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
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
	assert.Contains(t, resp.Body.String(), "token mismatch")
}
