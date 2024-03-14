package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexferl/echo-openapi"
	api "github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/services"
)

type UserHandlerTestSuite struct {
	suite.Suite
	svc              *handlers.MockUserService
	server           *api.Server
	user             *models.User
	accessToken      []byte
	admin            *models.User
	adminAccessToken []byte
	super            *models.User
}

func (s *UserHandlerTestSuite) SetupTest() {
	svc := handlers.NewMockUserService(s.T())
	patSvc := handlers.NewMockPersonalAccessTokenService(s.T())
	h := handlers.NewUserHandler(openapi.NewHandler(), svc)

	user := getUser()
	access, _, _ := user.Login()

	admin := getAdmin()
	adminAccess, _, _ := admin.Login()

	super := getSuper()

	s.svc = svc
	s.server = getServer(svc, patSvc, h)
	s.user = user
	s.accessToken = access
	s.admin = admin
	s.adminAccessToken = adminAccess
	s.super = super
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (s *UserHandlerTestSuite) TestUserHandler_GetCurrentUser_200() {
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result models.UserResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), s.user.Id, result.Id)
}

func (s *UserHandlerTestSuite) TestUserHandler_UpdateCurrentUser_200() {
	updatedUser := s.user
	updatedUser.Name = "updated name"
	b, _ := json.Marshal(&handlers.UpdateCurrentUserRequest{
		Name: &updatedUser.Name,
	})

	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(updatedUser, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result models.UserResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), updatedUser.Name, result.Name)
}

func (s *UserHandlerTestSuite) TestUserHandler_UpdateCurrentUser_422() {
	req := httptest.NewRequest(http.MethodPatch, "/me", bytes.NewBuffer([]byte(`{"invalid": "key"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.Code)
}

func (s *UserHandlerTestSuite) TestUserHandler_Get_200() {
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result models.UserResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), s.user.Id, result.Id)
}

func (s *UserHandlerTestSuite) TestUserHandler_Update_200() {
	updatedUser := s.user
	updatedUser.Name = "updated name"
	b, _ := json.Marshal(&handlers.UpdateUserRequest{
		Name: &updatedUser.Name,
	})

	req := httptest.NewRequest(http.MethodPatch, "/users/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.admin, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(updatedUser, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result models.UserResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), updatedUser.Name, result.Name)
}

func (s *UserHandlerTestSuite) TestUserHandler_Update_404() {
	updatedUser := s.user
	updatedUser.Name = "updated name"
	b, _ := json.Marshal(&handlers.UpdateUserRequest{
		Name: &updatedUser.Name,
	})

	req := httptest.NewRequest(http.MethodPatch, "/users/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.admin, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(nil, &services.Error{
			Kind:    services.NotExist,
			Message: services.ErrUserNotFound.Error(),
		}).Once()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
	assert.Equal(s.T(), services.ErrUserNotFound.Error(), result.Message)
}

func (s *UserHandlerTestSuite) TestUserHandler_Update_410() {
	updatedUser := s.user
	updatedUser.Name = "updated name"
	b, _ := json.Marshal(&handlers.UpdateUserRequest{
		Name: &updatedUser.Name,
	})
	updatedUser.Delete(s.admin.Id)

	req := httptest.NewRequest(http.MethodPatch, "/users/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.admin, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(nil, &services.Error{
			Kind:    services.Deleted,
			Message: services.ErrUserDeleted.Error(),
		}).Once()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusGone, resp.Code)
	assert.Equal(s.T(), services.ErrUserDeleted.Error(), result.Message)
}

func (s *UserHandlerTestSuite) TestUserHandler_204() {
	bannedUser := models.NewUser("banned@example.com", "banned")
	_ = bannedUser.Ban(s.admin)

	lockedUser := models.NewUser("locked@example.com", "locked")
	_ = lockedUser.Lock(s.admin)

	testCases := []struct {
		method   string
		endpoint string
		target   *models.User
	}{
		{http.MethodPut, "/users/1/ban", s.user},
		{http.MethodDelete, "/users/1/ban", bannedUser},
		{http.MethodPut, "/users/1/lock", s.user},
		{http.MethodDelete, "/users/1/lock", lockedUser},
		{http.MethodPut, "/users/1/roles/admin", s.user},
		{http.MethodDelete, "/users/1/roles/user", s.user},
	}
	for _, tc := range testCases {
		s.T().Run(fmt.Sprintf("%s_%s", tc.method, tc.endpoint), func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.endpoint, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
			resp := httptest.NewRecorder()

			// middleware
			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(s.admin, nil).Once()

			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(tc.target, nil).Once()

			s.svc.EXPECT().
				Update(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, nil).Once()

			s.server.ServeHTTP(resp, req)

			assert.Equal(s.T(), http.StatusNoContent, resp.Code)
		})
	}
}

func (s *UserHandlerTestSuite) TestUserHandler_401() {
	testCases := []struct {
		method   string
		endpoint string
	}{
		{http.MethodGet, "/me"},
		{http.MethodPatch, "/me"},
		{http.MethodGet, "/users/1"},
		{http.MethodPatch, "/users/1"},
		{http.MethodPut, "/users/1/ban"},
		{http.MethodDelete, "/users/1/ban"},
		{http.MethodPut, "/users/1/lock"},
		{http.MethodDelete, "/users/1/lock"},
		{http.MethodPut, "/users/1/roles/user"},
		{http.MethodDelete, "/users/1/roles/user"},
		{http.MethodGet, "/users"},
	}
	for _, tc := range testCases {
		s.T().Run(fmt.Sprintf("%s_%s", tc.method, tc.endpoint), func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.endpoint, nil)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			s.server.ServeHTTP(resp, req)

			var result echo.HTTPError
			_ = json.Unmarshal(resp.Body.Bytes(), &result)

			assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
			assert.Equal(s.T(), "token invalid", result.Message)
		})
	}
}

func (s *UserHandlerTestSuite) TestUserHandler_403() {
	testCases := []struct {
		method   string
		endpoint string
	}{
		{http.MethodPut, "/users/1/ban"},
		{http.MethodDelete, "/users/1/ban"},
		{http.MethodPut, "/users/1/lock"},
		{http.MethodDelete, "/users/1/lock"},
		{http.MethodPut, "/users/1/roles/super"},
		{http.MethodDelete, "/users/1/roles/super"},
	}
	for _, tc := range testCases {
		s.T().Run(fmt.Sprintf("%s_%s", tc.method, tc.endpoint), func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.endpoint, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
			resp := httptest.NewRecorder()

			// middleware
			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(s.admin, nil).Once()

			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(s.super, nil).Once()

			s.server.ServeHTTP(resp, req)

			var result echo.HTTPError
			_ = json.Unmarshal(resp.Body.Bytes(), &result)

			assert.Equal(s.T(), http.StatusForbidden, resp.Code)
		})
	}
}

func (s *UserHandlerTestSuite) TestUserHandler_404() {
	testCases := []struct {
		method   string
		endpoint string
	}{
		{http.MethodGet, "/users/1"},
		{http.MethodPut, "/users/1/ban"},
		{http.MethodDelete, "/users/1/ban"},
		{http.MethodPut, "/users/1/lock"},
		{http.MethodDelete, "/users/1/lock"},
		{http.MethodPut, "/users/1/roles/user"},
		{http.MethodDelete, "/users/1/roles/user"},
	}
	for _, tc := range testCases {
		s.T().Run(fmt.Sprintf("%s_%s", tc.method, tc.endpoint), func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.endpoint, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
			resp := httptest.NewRecorder()

			// middleware
			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(s.admin, nil).Once()

			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(nil, &services.Error{
					Kind:    services.NotExist,
					Message: services.ErrUserNotFound.Error(),
				}).Once()

			s.server.ServeHTTP(resp, req)

			var result echo.HTTPError
			_ = json.Unmarshal(resp.Body.Bytes(), &result)

			assert.Equal(s.T(), http.StatusNotFound, resp.Code)
			assert.Equal(s.T(), services.ErrUserNotFound.Error(), result.Message)
		})
	}
}

func (s *UserHandlerTestSuite) TestUserHandler_409() {
	testCases := []struct {
		method   string
		endpoint string
	}{
		{http.MethodPut, "/users/1/ban"},
		{http.MethodDelete, "/users/1/ban"},
		{http.MethodPut, "/users/1/lock"},
		{http.MethodDelete, "/users/1/lock"},
		{http.MethodPut, "/users/1/roles/user"},
		{http.MethodDelete, "/users/1/roles/user"},
	}
	for _, tc := range testCases {
		s.T().Run(fmt.Sprintf("%s_%s", tc.method, tc.endpoint), func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.endpoint, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
			resp := httptest.NewRecorder()

			// middleware
			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(s.admin, nil).Once()

			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(s.admin, nil).Once()

			s.server.ServeHTTP(resp, req)

			var result echo.HTTPError
			_ = json.Unmarshal(resp.Body.Bytes(), &result)

			assert.Equal(s.T(), http.StatusConflict, resp.Code)
		})
	}
}

func (s *UserHandlerTestSuite) TestUserHandler_410() {
	testCases := []struct {
		method   string
		endpoint string
	}{
		{http.MethodGet, "/users/1"},
		{http.MethodPut, "/users/1/ban"},
		{http.MethodDelete, "/users/1/ban"},
		{http.MethodPut, "/users/1/lock"},
		{http.MethodDelete, "/users/1/lock"},
		{http.MethodPut, "/users/1/roles/user"},
		{http.MethodDelete, "/users/1/roles/user"},
	}
	for _, tc := range testCases {
		s.T().Run(fmt.Sprintf("%s_%s", tc.method, tc.endpoint), func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.endpoint, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
			resp := httptest.NewRecorder()

			// middleware
			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(s.admin, nil).Once()

			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(nil, &services.Error{
					Kind:    services.Deleted,
					Message: services.ErrUserDeleted.Error(),
				}).Once()

			s.server.ServeHTTP(resp, req)

			var result echo.HTTPError
			_ = json.Unmarshal(resp.Body.Bytes(), &result)

			assert.Equal(s.T(), http.StatusGone, resp.Code)
			assert.Equal(s.T(), services.ErrUserDeleted.Error(), result.Message)
		})
	}
}

func (s *UserHandlerTestSuite) TestUserHandler_Roles_422() {
	testCases := []struct {
		method   string
		endpoint string
	}{
		{http.MethodPut, "/users/1/roles/wrong"},
		{http.MethodDelete, "/users/1/roles/wrong"},
	}
	for _, tc := range testCases {
		s.T().Run(fmt.Sprintf("%s_%s", tc.method, tc.endpoint), func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.endpoint, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
			resp := httptest.NewRecorder()

			// middleware
			s.svc.EXPECT().
				Read(mock.Anything, mock.Anything).
				Return(s.admin, nil).Once()

			s.server.ServeHTTP(resp, req)

			var result echo.HTTPError
			_ = json.Unmarshal(resp.Body.Bytes(), &result)

			assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.Code)
		})
	}
}

func createUsers(num int) models.Users {
	result := make(models.Users, 0)

	for i := range num {
		user := models.NewUser(fmt.Sprintf("user%d@example.com", i), fmt.Sprintf("user%d", i))
		result = append(result, *user)
	}

	return result
}

func (s *UserHandlerTestSuite) TestUserHandler_List_200() {
	req := httptest.NewRequest(http.MethodGet, "/users?per_page=1&page=2", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
	resp := httptest.NewRecorder()

	num := 10
	users := createUsers(num)

	// middleware
	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.admin, nil).Once()

	s.svc.EXPECT().
		Find(mock.Anything, mock.Anything).
		Return(int64(num), users, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result models.PublicUsersResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	h := resp.Header()
	link := `<http://example.com/users?per_page=1&page=3>; rel=next, ` +
		`<http://example.com/users?per_page=1&page=10>; rel=last, ` +
		`<http://example.com/users?per_page=1&page=1>; rel=first, ` +
		`<http://example.com/users?per_page=1&page=1>; rel=prev`

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), 10, len(result.Users))
	assert.Equal(s.T(), "2", h.Get("X-Page"))
	assert.Equal(s.T(), "1", h.Get("X-Per-Page"))
	assert.Equal(s.T(), "10", h.Get("X-Total"))
	assert.Equal(s.T(), "10", h.Get("X-Total-Pages"))
	assert.Equal(s.T(), "3", h.Get("X-Next-Page"))
	assert.Equal(s.T(), "1", h.Get("X-Prev-Page"))
	assert.Equal(s.T(), link, h.Get("Link"))
}

func (s *UserHandlerTestSuite) TestUserHandler_List_403() {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusForbidden, resp.Code)
}
