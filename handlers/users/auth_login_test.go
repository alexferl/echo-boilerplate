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
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers/users"
)

func TestHandler_AuthLogin_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	pwd := "abcdefghijkl"
	user := users.NewUser("test@example.com", "test")
	err := user.SetPassword(pwd)
	assert.NoError(t, err)

	b, err := json.Marshal(&users.AuthLogInRequest{Password: pwd})
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
			"UpdateOneById",
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

func TestHandler_AuthLogin_401(t *testing.T) {
	payload := &users.AuthLogInRequest{}
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
			data.ErrNoDocuments,
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

func TestHandler_AuthLogin_422(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer([]byte(`{"invalid": "key"}`)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}
