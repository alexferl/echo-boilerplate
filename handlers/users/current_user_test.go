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
)

func TestHandler_GetCurrentUser_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
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

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHandler_GetCurrentUser_401(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHandler_UpdateCurentUser_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	user.Name = "test name"
	user.Bio = "test bio"
	access, _, err := user.Login()
	assert.NoError(t, err)

	updatedUser := user
	updatedUser.Name = "name"
	updatedUser.Bio = "bio"
	updatedUser.Update(user.Id)

	b, err := json.Marshal(&users.UpdateCurrentUserRequest{
		Name: updatedUser.Name,
		Bio:  updatedUser.Bio,
	})
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/me", bytes.NewBuffer(b))
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
		On("FindOneByIdAndUpdate",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			updatedUser,
			nil,
		)

	s.ServeHTTP(resp, req)

	fmt.Println(resp)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHandler_UpdateCurrentUser_401(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPut, "/me", bytes.NewBuffer([]byte(`{"invalid": "key"}`)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHandler_UpdateCurrentUser_422(t *testing.T) {
	_, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	user.Name = "test name"
	user.Bio = "test bio"
	access, _, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/me", bytes.NewBuffer([]byte(`{"invalid": "key"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}
