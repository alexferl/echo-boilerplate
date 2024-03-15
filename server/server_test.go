package server

import (
	"encoding/json"
	"errors"
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
	_ "github.com/alexferl/echo-boilerplate/testing"
	"github.com/alexferl/echo-boilerplate/util/cookie"
)

type ServerTestSuite struct {
	suite.Suite
	svc         *handlers.MockUserService
	patSvc      *handlers.MockPersonalAccessTokenService
	server      *api.Server
	user        *models.User
	accessToken []byte
	admin       *models.User
}

func (s *ServerTestSuite) SetupTest() {
	svc := handlers.NewMockUserService(s.T())
	patSvc := handlers.NewMockPersonalAccessTokenService(s.T())
	h := handlers.NewUserHandler(openapi.NewHandler(), svc)

	admin := models.NewUserWithRole("test@example.com", "test", models.AdminRole)
	user := models.NewUser("test@example.com", "test")
	user.Id = "1000"
	user.Create(user.Id)
	access, _, _ := user.Login()

	s.svc = svc
	s.patSvc = patSvc
	s.server = NewTestServer(svc, patSvc, h)
	s.user = user
	s.accessToken = access
	s.admin = admin
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) TestServer_503() {
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(nil, errors.New("")).Once()

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusServiceUnavailable, resp.Code)
}

func (s *ServerTestSuite) TestServer_403_Banned() {
	_ = s.user.Ban(s.admin)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusForbidden, resp.Code)
	s.Assert().Equal(ErrBanned.Error(), result.Message)
}

func (s *ServerTestSuite) TestServer_403_Locked() {
	_ = s.user.Lock(s.admin)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusForbidden, resp.Code)
	s.Assert().Equal(ErrLocked.Error(), result.Message)
}

func (s *ServerTestSuite) TestServer_400_CSRF_Header_Missing() {
	req := httptest.NewRequest(http.MethodPatch, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie.NewAccessToken(s.accessToken))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusBadRequest, resp.Code)
	s.Assert().Equal(ErrCSRFHeaderMissing.Error(), result.Message)
}

func (s *ServerTestSuite) TestServer_400_CSRF_Header_Invalid() {
	req := httptest.NewRequest(http.MethodPatch, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie.NewAccessToken(s.accessToken))
	req.Header.Add("X-CSRF-Token", "token")
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusForbidden, resp.Code)
	s.Assert().Equal(ErrCSRFInvalid.Error(), result.Message)
}

func (s *ServerTestSuite) TestServer_PAT_401_Token_Invalid() {
	pat, _ := models.NewPersonalAccessToken(
		s.user.Id,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat.Token))
	resp := httptest.NewRecorder()

	_ = pat.Encrypt()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.patSvc.EXPECT().
		FindOne(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, services.NewError(nil, services.NotExist, "")).Once()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
	s.Assert().Equal(ErrTokenInvalid.Error(), result.Message)
}

func (s *ServerTestSuite) TestServer_PAT_401_Token_Mismatch() {
	pat, _ := models.NewPersonalAccessToken(
		s.user.Id,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat.Token))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.patSvc.EXPECT().
		FindOne(mock.Anything, mock.Anything, mock.Anything).
		Return(pat, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
	s.Assert().Equal(ErrTokenMismatch.Error(), result.Message)
}

func (s *ServerTestSuite) TestServer_PAT_401_Revoked() {
	pat, _ := models.NewPersonalAccessToken(
		s.user.Id,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat.Token))
	resp := httptest.NewRecorder()

	_ = pat.Encrypt()
	pat.IsRevoked = true

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.patSvc.EXPECT().
		FindOne(mock.Anything, mock.Anything, mock.Anything).
		Return(pat, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
	s.Assert().Equal(ErrTokenRevoked.Error(), result.Message)
}

func (s *ServerTestSuite) TestServer_PAT_401_Expired() {
	pat, _ := models.NewPersonalAccessToken(
		s.user.Id,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)

	past := time.Now().Add(-(7 * 24) * time.Hour)
	pat.ExpiresAt = &past

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat.Token))
	resp := httptest.NewRecorder()

	_ = pat.Encrypt()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.patSvc.EXPECT().
		FindOne(mock.Anything, mock.Anything, mock.Anything).
		Return(pat, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
	s.Assert().Equal(ErrTokenExpired.Error(), result.Message)
}

func (s *ServerTestSuite) TestServer_PAT_503() {
	pat, _ := models.NewPersonalAccessToken(
		s.user.Id,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat.Token))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.patSvc.EXPECT().
		FindOne(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("")).Once()

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusServiceUnavailable, resp.Code)
}

func (s *ServerTestSuite) TestServer_PAT_200() {
	pat, _ := models.NewPersonalAccessToken(
		s.user.Id,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pat.Token))
	resp := httptest.NewRecorder()

	_ = pat.Encrypt()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.patSvc.EXPECT().
		FindOne(mock.Anything, mock.Anything, mock.Anything).
		Return(pat, nil).Once()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(s.user, nil).Once()

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusOK, resp.Code)
}
