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
	jwx "github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/server"
	"github.com/alexferl/echo-boilerplate/services"
	"github.com/alexferl/echo-boilerplate/util/jwt"
)

type PersonalAccessTokenHandlerTestSuite struct {
	suite.Suite
	svc              *handlers.MockPersonalAccessTokenService
	server           *api.Server
	user             *models.User
	accessToken      []byte
	admin            *models.User
	adminAccessToken []byte
}

func (s *PersonalAccessTokenHandlerTestSuite) SetupTest() {
	svc := handlers.NewMockPersonalAccessTokenService(s.T())
	h := handlers.NewPersonalAccessTokenHandler(openapi.NewHandler(), svc)
	user := models.NewUser("test@example.com", "test")
	user.Id = "100"
	user.Create(user.Id)
	access, _, _ := user.Login()

	s.svc = svc
	s.server = server.NewTestServer(h)
	s.user = user
	s.accessToken = access
}

func TestPersonalAccessTokenHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PersonalAccessTokenHandlerTestSuite))
}

func createTokens(token jwx.Token, num int) models.PersonalAccessTokens {
	result := make(models.PersonalAccessTokens, 0)

	for i := range num {
		pat, _ := models.NewPersonalAccessToken(
			token,
			fmt.Sprintf("my_token%d", i),
			time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
		)
		result = append(result, *pat)
	}

	return result
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Create_200() {
	token, _ := jwt.ParseEncoded(s.accessToken)

	payload := &handlers.CreatePersonalAccessTokenRequest{
		Name:      "My Token",
		ExpiresAt: time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02"),
	}
	b, _ := json.Marshal(payload)

	newPAT, _ := models.NewPersonalAccessToken(token, payload.Name, payload.ExpiresAt)

	req := httptest.NewRequest(http.MethodPost, "/me/personal_access_tokens", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		FindOne(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)

	s.svc.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(newPAT, nil)

	s.server.ServeHTTP(resp, req)

	var result models.PersonalAccessToken
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Create_401() {
	req := httptest.NewRequest(http.MethodPost, "/me/personal_access_tokens", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Create_409() {
	token, _ := jwt.ParseEncoded(s.accessToken)

	payload := &handlers.CreatePersonalAccessTokenRequest{
		Name:      "My Token",
		ExpiresAt: time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02"),
	}
	b, _ := json.Marshal(payload)

	newPAT, _ := models.NewPersonalAccessToken(token, payload.Name, payload.ExpiresAt)

	req := httptest.NewRequest(http.MethodPost, "/me/personal_access_tokens", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		FindOne(mock.Anything, mock.Anything, mock.Anything).
		Return(newPAT, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusConflict, resp.Code)
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

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Create_422_Exp() {
	token, _ := jwt.ParseEncoded(s.accessToken)

	payload := &handlers.CreatePersonalAccessTokenRequest{
		Name:      "My Token",
		ExpiresAt: time.Now().Add(-(7 * 24) * time.Hour).Format("2006-01-02"),
	}
	b, _ := json.Marshal(payload)

	newPAT, _ := models.NewPersonalAccessToken(token, payload.Name, payload.ExpiresAt)

	req := httptest.NewRequest(http.MethodPost, "/me/personal_access_tokens", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		FindOne(mock.Anything, mock.Anything, mock.Anything).
		Return(newPAT, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_List_200() {
	token, _ := jwt.ParseEncoded(s.accessToken)
	num := 10
	pats := createTokens(token, num)

	req := httptest.NewRequest(http.MethodGet, "/me/personal_access_tokens", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Find(mock.Anything, mock.Anything).
		Return(pats, nil)

	s.server.ServeHTTP(resp, req)

	var result models.PersonalAccessTokensResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), num, len(result.Tokens))
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_List_401() {
	req := httptest.NewRequest(http.MethodGet, "/me/personal_access_tokens", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Get_200() {
	token, _ := jwt.ParseEncoded(s.accessToken)

	newPAT, _ := models.NewPersonalAccessToken(
		token,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)

	req := httptest.NewRequest(http.MethodGet, "/me/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(newPAT, nil)

	s.server.ServeHTTP(resp, req)

	var result models.PersonalAccessToken
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Get_404() {
	req := httptest.NewRequest(http.MethodGet, "/me/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(nil, &services.Error{
			Kind:    services.NotExist,
			Message: services.ErrPersonalAccessTokenNotFound.Error(),
		})

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
	assert.Equal(s.T(), services.ErrPersonalAccessTokenNotFound.Error(), result.Message)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Revoke_204() {
	token, _ := jwt.ParseEncoded(s.accessToken)

	newPAT, _ := models.NewPersonalAccessToken(
		token,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)

	req := httptest.NewRequest(http.MethodDelete, "/me/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(newPAT, nil)

	s.svc.EXPECT().
		Revoke(mock.Anything, mock.Anything).
		Return(nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNoContent, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Revoke_401() {
	req := httptest.NewRequest(http.MethodDelete, "/me/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Revoke_404() {
	req := httptest.NewRequest(http.MethodDelete, "/me/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(nil, &services.Error{
			Kind:    services.NotExist,
			Message: services.ErrPersonalAccessTokenNotFound.Error(),
		})

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
	assert.Equal(s.T(), services.ErrPersonalAccessTokenNotFound.Error(), result.Message)
}

func (s *PersonalAccessTokenHandlerTestSuite) TestPersonalAccessTokenHandler_Revoke_409() {
	token, _ := jwt.ParseEncoded(s.accessToken)

	newPAT, _ := models.NewPersonalAccessToken(
		token,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)
	newPAT.IsRevoked = true

	req := httptest.NewRequest(http.MethodDelete, "/me/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(newPAT, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusConflict, resp.Code)
}
