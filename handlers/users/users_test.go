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

func createUsers(num int) []*users.ShortUser {
	var result []*users.ShortUser

	for i := 1; i <= num; i++ {
		user := users.NewUser(fmt.Sprintf("user%d@example.com", i), fmt.Sprintf("user%d", i))
		short := &users.ShortUser{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
		}
		result = append(result, short)
	}

	return result
}

func TestHandler_Users_200(t *testing.T) {
	retUsers := createUsers(10)
	testCases := []struct {
		query       string
		link        string
		xPage       string
		xPerPage    string
		xTotal      string
		xTotalPages string
		xNextPage   string
		xPrevPage   string
	}{
		{
			query: "",
			link: `<http://example.com/users?per_page=10&page=1>; rel=last, ` +
				`<http://example.com/users?per_page=10&page=1>; rel=first`,
			xPage:       "1",
			xPerPage:    "10",
			xTotal:      "10",
			xTotalPages: "1",
			xNextPage:   "",
			xPrevPage:   "",
		},
		{
			query: "per_page=10&page=1",
			link: `<http://example.com/users?per_page=10&page=1>; rel=last, ` +
				`<http://example.com/users?per_page=10&page=1>; rel=first`,
			xPage:       "1",
			xPerPage:    "10",
			xTotal:      "10",
			xTotalPages: "1",
			xNextPage:   "",
			xPrevPage:   "",
		},
		{
			query: "per_page=1&page=2",
			link: `<http://example.com/users?per_page=1&page=3>; rel=next, ` +
				`<http://example.com/users?per_page=1&page=10>; rel=last, ` +
				`<http://example.com/users?per_page=1&page=1>; rel=first, ` +
				`<http://example.com/users?per_page=1&page=1>; rel=prev`,
			xPage:       "2",
			xPerPage:    "1",
			xTotal:      "10",
			xTotalPages: "10",
			xNextPage:   "3",
			xPrevPage:   "1",
		},
		{
			query: "per_page=1&page=10",
			link: `<http://example.com/users?per_page=1&page=10>; rel=last, ` +
				`<http://example.com/users?per_page=1&page=1>; rel=first, ` +
				`<http://example.com/users?per_page=1&page=9>; rel=prev`,
			xPage:       "10",
			xPerPage:    "1",
			xTotal:      "10",
			xTotalPages: "10",
			xNextPage:   "",
			xPrevPage:   "9",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.query, func(t *testing.T) {
			mapper, s := getMapperAndServer(t)

			user := users.NewAdminUser("admin@example.com", "admin")
			access, _, err := user.Login()
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/users?"+tc.query, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
			resp := httptest.NewRecorder()

			mapper.Mock.
				On(
					"Count",
					mock.Anything,
					mock.Anything,
				).
				Return(
					int64(10),
					nil,
				).
				On(
					"Find",
					mock.Anything,
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

			h := resp.Header()

			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, 10, len(result.Users))
			assert.Equal(t, tc.xPage, h.Get("X-Page"))
			assert.Equal(t, tc.xPerPage, h.Get("X-Per-Page"))
			assert.Equal(t, tc.xTotal, h.Get("X-Total"))
			assert.Equal(t, tc.xTotalPages, h.Get("X-Total-Pages"))
			assert.Equal(t, tc.xNextPage, h.Get("X-Next-Page"))
			assert.Equal(t, tc.xPrevPage, h.Get("X-Prev-Page"))
			assert.Equal(t, tc.link, h.Get("Link"))
		})
	}
}

func TestHandler_Users_401(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
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
