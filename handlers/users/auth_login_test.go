package users_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Auth_Login_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	pwd := "abcdefghijkl"
	user := users.NewUser("test@example.com", "test")
	err := user.SetPassword(pwd)
	assert.NoError(t, err)

	b, err := json.Marshal(&users.LoginPayload{Password: pwd})
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(b))
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

func TestHandler_Auth_Login_401(t *testing.T) {
	payload := &users.LoginPayload{}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	testCases := []struct {
		name       string
		payload    []byte
		returnUser *users.User
		err        error
		statusCode int
		msg        string
	}{
		{
			"user not found",
			b,
			nil,
			users.ErrUserNotFound,
			http.StatusUnauthorized,
			"invalid email or password",
		},
		{
			"wrong password",
			b,
			&users.User{},
			nil,
			http.StatusUnauthorized,
			"invalid email or password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mapper, s := getMapperAndServer(t)

			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(tc.payload))
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
					tc.returnUser,
					tc.err,
				)
			s.ServeHTTP(resp, req)
			assert.Equal(t, tc.statusCode, resp.Code)
			if tc.msg != "" {
				assert.Contains(t, resp.Body.String(), tc.msg)
			}
		})
	}
}

func TestHandler_Auth_Login_422(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}
