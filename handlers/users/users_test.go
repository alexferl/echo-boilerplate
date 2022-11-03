package users_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alexferl/echo-boilerplate/handlers/users"
)

func TestHandler_Users_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewAdminUser("admin@example.com", "admin")
	access, _, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	user1 := users.NewUser("user1@example.com", "user1")
	short1 := &users.ShortUser{
		Id:       user1.Id,
		Username: user1.Username,
		Email:    user1.Email,
	}

	user2 := users.NewUser("user2@example.com", "user2")
	short2 := &users.ShortUser{
		Id:       user2.Id,
		Username: user2.Username,
		Email:    user2.Email,
	}

	retUsers := []*users.ShortUser{short1, short2}

	mapper.Mock.
		On(
			"Find",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			retUsers,
			nil,
		)

	s.ServeHTTP(resp, req)

	var result users.UsersResponse
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, 2, len(result.Users))
}

func TestHandler_Users_403(t *testing.T) {
	_, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusForbidden, resp.Code)
}
