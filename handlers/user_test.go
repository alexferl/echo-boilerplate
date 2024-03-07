package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	app "github.com/alexferl/echo-boilerplate"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/models"
)

type UserHandlerTestSuite struct {
	suite.Suite
	svc              *handlers.MockUserService
	server           *server.Server
	user             *models.User
	accessToken      []byte
	admin            *models.User
	adminAccessToken []byte
}

func (s *UserHandlerTestSuite) SetupTest() {
	svc := handlers.NewMockUserService(s.T())
	h := handlers.NewUserHandler(openapi.NewHandler(), svc)
	user := models.NewUser("test@example.com", "test")
	user.Id = "100"
	user.Create(user.Id)
	access, _, _ := user.Login()

	admin := models.NewUserWithRole("admin@example.com", "admin", models.AdminRole)
	admin.Id = "200"
	admin.Create(admin.Id)
	adminAccess, _, _ := admin.Login()

	s.svc = svc
	s.server = app.NewTestServer(h)
	s.user = user
	s.accessToken = access
	s.admin = admin
	s.adminAccessToken = adminAccess
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (s *UserHandlerTestSuite) TestUserHandler_GetCurrentUser_200() {
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil)

	s.server.ServeHTTP(resp, req)

	var result models.UserResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), s.user.Id, result.Id)
}

func (s *UserHandlerTestSuite) TestUserHandler_GetCurrentUser_401() {
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), "token invalid", result.Message)
}

func (s *UserHandlerTestSuite) TestUserHandler_UpdateCurrentUser_200() {
	updatedUser := s.user
	updatedUser.Name = "updated name"
	b, _ := json.Marshal(&handlers.UpdateCurrentUserRequest{
		Name: &updatedUser.Name,
	})

	req := httptest.NewRequest(http.MethodPut, "/me", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil)

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(updatedUser, nil)

	s.server.ServeHTTP(resp, req)

	var result models.UserResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), updatedUser.Name, result.Name)
}

func (s *UserHandlerTestSuite) TestUserHandler_UpdateCurrentUser_401() {
	req := httptest.NewRequest(http.MethodPut, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *UserHandlerTestSuite) TestUserHandler_UpdateCurrentUser_422() {
	req := httptest.NewRequest(http.MethodPut, "/me", bytes.NewBuffer([]byte(`{"invalid": "key"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.Code)
}

func (s *UserHandlerTestSuite) TestUserHandler_Get_200() {
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil)

	s.server.ServeHTTP(resp, req)

	var result models.UserResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), s.user.Id, result.Id)
}

func (s *UserHandlerTestSuite) TestUserHandler_Get_404() {
	req := httptest.NewRequest(http.MethodGet, "/users/404", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
}

func (s *UserHandlerTestSuite) TestUserHandler_Get_410() {
	user := models.NewUser("deleted@example.com", "deleted")
	user.Delete(user.Id)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusGone, resp.Code)
}

func (s *UserHandlerTestSuite) TestUserHandler_Update_200() {
	updatedUser := s.user
	updatedUser.Name = "updated name"
	b, _ := json.Marshal(&handlers.UpdateUserRequest{
		Name: &updatedUser.Name,
	})

	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil)

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(updatedUser, nil)

	s.server.ServeHTTP(resp, req)

	var result models.UserResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), updatedUser.Name, result.Name)
}

func (s *UserHandlerTestSuite) TestUserHandler_Update_401() {
	req := httptest.NewRequest(http.MethodPut, "/users/1", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *UserHandlerTestSuite) TestUserHandler_Update_404() {
	updatedUser := s.user
	updatedUser.Name = "updated name"
	b, _ := json.Marshal(&handlers.UpdateUserRequest{
		Name: &updatedUser.Name,
	})

	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
}

func (s *UserHandlerTestSuite) TestUserHandler_Update_410() {
	updatedUser := s.user
	updatedUser.Name = "updated name"
	b, _ := json.Marshal(&handlers.UpdateUserRequest{
		Name: &updatedUser.Name,
	})
	updatedUser.Delete(s.admin.Id)

	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(updatedUser, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusGone, resp.Code)
}

func (s *UserHandlerTestSuite) TestUserHandler_UpdateStatus_200() {
	updatedUser := s.user
	updatedUser.Name = "updated name"
	t := true
	b, _ := json.Marshal(&handlers.UpdateUserStatusRequest{
		IsLocked: &t,
	})

	req := httptest.NewRequest(http.MethodPut, "/users/1/status", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil)

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(updatedUser, nil)

	s.server.ServeHTTP(resp, req)

	var result models.UserResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), updatedUser.Name, result.Name)
}

func (s *UserHandlerTestSuite) TestUserHandler_UpdateStatus_401() {
	req := httptest.NewRequest(http.MethodPut, "/users/1/status", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *UserHandlerTestSuite) TestUserHandler_UpdateStatus_404() {
	updatedUser := s.user
	updatedUser.Name = "updated name"
	t := true
	b, _ := json.Marshal(&handlers.UpdateUserStatusRequest{
		IsLocked: &t,
	})

	req := httptest.NewRequest(http.MethodPut, "/users/1/status", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
}

func (s *UserHandlerTestSuite) TestUserHandler_UpdateStatus_410() {
	updatedUser := s.user
	updatedUser.Name = "updated name"
	t := true
	b, _ := json.Marshal(&handlers.UpdateUserStatusRequest{
		IsLocked: &t,
	})
	updatedUser.Delete(s.admin.Id)

	req := httptest.NewRequest(http.MethodPut, "/users/1/status", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(updatedUser, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusGone, resp.Code)
}

func (s *UserHandlerTestSuite) TestUserHandler_UpdateStatus_409() {
	updatedUser := s.admin
	updatedUser.Name = "updated name"
	t := true
	b, _ := json.Marshal(&handlers.UpdateUserStatusRequest{
		IsLocked: &t,
	})

	req := httptest.NewRequest(http.MethodPut, "/users/1/status", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.adminAccessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(updatedUser, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusConflict, resp.Code)
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

	s.svc.EXPECT().
		Find(mock.Anything, mock.Anything).
		Return(int64(num), users, nil)

	s.server.ServeHTTP(resp, req)

	var result handlers.ListUsersResponse
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

func (s *UserHandlerTestSuite) TestUserHandler_List_401() {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *UserHandlerTestSuite) TestUserHandler_List_403() {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusForbidden, resp.Code)
}
