package users_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alexferl/echo-boilerplate/handlers/users"
)

func TestHandler_GetUsername_200(t *testing.T) {
	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	result := &users.GetUsernameResponse{
		Id:        user.Id,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		DeletedAt: user.DeletedAt,
	}

	testCases := []struct {
		name       string
		username   string
		statusCode int
		retUser    *users.GetUsernameResponse
		retErr     error
	}{
		{
			"not found", "notfound", http.StatusNotFound, nil, users.ErrUserNotFound,
		},
		{
			"self username", "test", http.StatusOK, result, nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mapper, s := getMapperAndServer(t)

			target := fmt.Sprintf("/users/%s", tc.username)
			req := httptest.NewRequest(http.MethodGet, target, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
			resp := httptest.NewRecorder()

			mapper.Mock.
				On(
					"FindOne",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
				Return(
					tc.retUser,
					tc.retErr,
				)

			s.ServeHTTP(resp, req)

			assert.Equal(t, tc.statusCode, resp.Code)
		})
	}
}

func TestHandler_GetUsername_200_Not_Logged_In(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")

	result := &users.GetUsernameResponse{
		Id:        user.Id,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		DeletedAt: user.DeletedAt,
	}

	req := httptest.NewRequest(http.MethodGet, "/users/test", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOne",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			result,
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHandler_GetUsername_404(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/users/notfound", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOne",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			nil,
			users.ErrUserNotFound,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestHandler_GetUsername_410(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	user.Delete(user.Id)
	access, _, err := user.Login()
	assert.NoError(t, err)

	result := &users.GetUsernameResponse{
		Id:        user.Id,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		DeletedAt: user.DeletedAt,
	}

	req := httptest.NewRequest(http.MethodGet, "/users/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"FindOne",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			result,
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusGone, resp.Code)
}
