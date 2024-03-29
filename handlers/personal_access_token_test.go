package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexferl/echo-openapi"
	api "github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/services"
)

type PersonalAccessTokenHandlerTestSuite struct {
	suite.Suite
	svc              *handlers.MockPersonalAccessTokenService
	userSvc          *handlers.MockUserService
	server           *api.Server
	user             *models.User
	accessToken      []byte
	admin            *models.User
	adminAccessToken []byte
}

func (s *PersonalAccessTokenHandlerTestSuite) SetupTest() {
	userSvc := handlers.NewMockUserService(s.T())
	svc := handlers.NewMockPersonalAccessTokenService(s.T())
	h := handlers.NewPersonalAccessTokenHandler(openapi.NewHandler(), svc)
	user := getUser()
	access, _, _ := user.Login()

	s.svc = svc
	s.userSvc = userSvc
	s.server = getServer(userSvc, svc, h)
	s.user = user
	s.accessToken = access
}

func TestPersonalAccessTokenHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PersonalAccessTokenHandlerTestSuite))
}

func createTokens(userId string, num int) models.PersonalAccessTokens {
	result := make(models.PersonalAccessTokens, 0)

	for i := range num {
		pat, _ := models.NewPersonalAccessToken(
			userId,
			fmt.Sprintf("my_token%d", i),
			time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
		)
		result = append(result, *pat)
	}

	return result
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Create_200() {
	payload := &handlers.CreatePersonalAccessTokenRequest{
		Name:      "My Token",
		ExpiresAt: time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02"),
	}
	b, _ := json.Marshal(payload)

	newPAT, _ := models.NewPersonalAccessToken(s.user.Id, payload.Name, payload.ExpiresAt)

	req := httptest.NewRequest(http.MethodPost, "/me/personal_access_tokens", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		FindOne(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)

	s.svc.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(newPAT, nil)

	s.server.ServeHTTP(resp, req)

	var result models.PersonalAccessToken
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusOK, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Create_401() {
	req := httptest.NewRequest(http.MethodPost, "/me/personal_access_tokens", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Create_409() {
	payload := &handlers.CreatePersonalAccessTokenRequest{
		Name:      "My Token",
		ExpiresAt: time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02"),
	}
	b, _ := json.Marshal(payload)

	newPAT, _ := models.NewPersonalAccessToken(s.user.Id, payload.Name, payload.ExpiresAt)

	req := httptest.NewRequest(http.MethodPost, "/me/personal_access_tokens", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		FindOne(mock.Anything, mock.Anything, mock.Anything).
		Return(newPAT, nil)

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusConflict, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Create_422() {
	payload := &handlers.CreatePersonalAccessTokenRequest{
		Name:      "",
		ExpiresAt: "invalid",
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/me/personal_access_tokens", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusUnprocessableEntity, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Create_422_Exp() {
	payload := &handlers.CreatePersonalAccessTokenRequest{
		Name:      "My Token",
		ExpiresAt: time.Now().Add(-(7 * 24) * time.Hour).Format("2006-01-02"),
	}
	b, _ := json.Marshal(payload)

	newPAT, _ := models.NewPersonalAccessToken(s.user.Id, payload.Name, payload.ExpiresAt)

	req := httptest.NewRequest(http.MethodPost, "/me/personal_access_tokens", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		FindOne(mock.Anything, mock.Anything, mock.Anything).
		Return(newPAT, nil)

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusUnprocessableEntity, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_List_200() {
	num := 10
	pats := createTokens(s.user.Id, num)

	req := httptest.NewRequest(http.MethodGet, "/me/personal_access_tokens", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Find(mock.Anything, mock.Anything).
		Return(pats, nil)

	s.server.ServeHTTP(resp, req)

	var result models.PersonalAccessTokensResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().Equal(num, len(result.Tokens))
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_List_401() {
	req := httptest.NewRequest(http.MethodGet, "/me/personal_access_tokens", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Get_200() {
	newPAT, _ := models.NewPersonalAccessToken(
		s.user.Id,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)

	req := httptest.NewRequest(http.MethodGet, "/me/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything, mock.Anything).
		Return(newPAT, nil)

	s.server.ServeHTTP(resp, req)

	var result models.PersonalAccessToken
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusOK, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Get_404() {
	req := httptest.NewRequest(http.MethodGet, "/me/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, &services.Error{
			Kind:    services.NotExist,
			Message: services.ErrPersonalAccessTokenNotFound.Error(),
		})

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusNotFound, resp.Code)
	s.Assert().Equal(services.ErrPersonalAccessTokenNotFound.Error(), result.Message)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Revoke_204() {
	newPAT, _ := models.NewPersonalAccessToken(
		s.user.Id,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)

	req := httptest.NewRequest(http.MethodDelete, "/me/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything, mock.Anything).
		Return(newPAT, nil)

	s.svc.EXPECT().
		Revoke(mock.Anything, mock.Anything).
		Return(nil)

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusNoContent, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Revoke_401() {
	req := httptest.NewRequest(http.MethodDelete, "/me/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Revoke_404() {
	req := httptest.NewRequest(http.MethodDelete, "/me/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, &services.Error{
			Kind:    services.NotExist,
			Message: services.ErrPersonalAccessTokenNotFound.Error(),
		})

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusNotFound, resp.Code)
	s.Assert().Equal(services.ErrPersonalAccessTokenNotFound.Error(), result.Message)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Revoke_409() {
	newPAT, _ := models.NewPersonalAccessToken(
		s.user.Id,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)
	newPAT.IsRevoked = true

	req := httptest.NewRequest(http.MethodDelete, "/me/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	// middleware
	s.userSvc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything, mock.Anything).
		Return(newPAT, nil)

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusConflict, resp.Code)
}
