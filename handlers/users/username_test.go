package users_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers/users"
)

func TestHandler_GetUsername_200(t *testing.T) {
	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	testCases := []struct {
		name       string
		username   string
		statusCode int
		retUser    *users.User
		retErr     error
	}{
		{
			"not found", "notfound", http.StatusNotFound, nil, data.ErrNoDocuments,
		},
		{
			"self username", "test", http.StatusOK, user, nil,
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
			user,
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
			data.ErrNoDocuments,
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
			user,
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusGone, resp.Code)
}
