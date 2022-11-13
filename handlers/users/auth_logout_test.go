package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
)

func TestHandler_AuthLogout_400_Cookie_Missing(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Request malformed")
}

func TestHandler_AuthLogout_401_Cookie_Invalid(t *testing.T) {
	_, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	_, _, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewRefreshTokenCookie([]byte("invalid")))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Token invalid")
}

func TestHandler_AuthLogout_401_Cookie_Mismatch(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	_, refresh, err := user.Login()
	assert.NoError(t, err)
	user.Logout()

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
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

func TestHandler_AuthLogout_200_Token(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	_, refresh, err := user.Login()
	assert.NoError(t, err)

	payload := &users.AuthLogOutRequest{
		RefreshToken: string(refresh),
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBuffer(b))
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
	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestHandler_AuthLogout_400_Token_Missing(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Request malformed")
}

func TestHandler_AuthLogout_401_Token_Invalid(t *testing.T) {
	_, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &users.AuthLogOutRequest{
		RefreshToken: "invalid",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "Token invalid")
}

func TestHandler_AuthLogout_401_Token_Mismatch(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	_, refresh, err := user.Login()
	assert.NoError(t, err)
	user.Logout()

	payload := &users.AuthLogOutRequest{
		RefreshToken: string(refresh),
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBuffer(b))
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
