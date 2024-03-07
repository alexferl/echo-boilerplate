package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexferl/echo-jwt"
	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	app "github.com/alexferl/echo-boilerplate"
	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/util"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	svc    *handlers.MockUserService
	server *server.Server
}

func (s *AuthHandlerTestSuite) SetupTest() {
	svc := handlers.NewMockUserService(s.T())
	h := handlers.NewAuthHandler(openapi.NewHandler(), svc)
	s.svc = svc
	s.server = app.NewTestServer(h)
}

func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Login_200() {
	pwd := "abcdefghijkl"
	user := models.NewUser("test@example.com", "test")
	_ = user.SetPassword(pwd)

	email, _ := json.Marshal(&handlers.LoginRequest{Email: user.Email, Password: pwd})
	username, _ := json.Marshal(&handlers.LoginRequest{Username: user.Username, Password: pwd})

	testCases := []struct {
		name    string
		payload []byte
	}{
		{"using email", email},
		{"using username", username},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			s.svc.EXPECT().
				FindOneByEmailOrUsername(mock.Anything, mock.Anything, mock.Anything).
				Return(user, nil)

			s.svc.EXPECT().
				Update(mock.Anything, mock.Anything, mock.Anything).
				Return(user, nil)

			s.server.ServeHTTP(resp, req)

			var result handlers.LoginResponse
			_ = json.Unmarshal(resp.Body.Bytes(), &result)

			expected := 2
			if viper.GetBool(config.CSRFEnabled) {
				expected = 3
			}
			if assert.Equal(s.T(), expected, len(resp.Result().Cookies())) {
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
				assert.Equal(s.T(), expected, cookies)
			}

			assert.Equal(s.T(), http.StatusOK, resp.Code)
			assert.NotEqual(s.T(), "", result.AccessToken)
			assert.NotEqual(s.T(), "", result.ExpiresIn)
			assert.NotEqual(s.T(), "", result.RefreshToken)
			assert.NotEqual(s.T(), "", result.TokenType)
		})
	}
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Login_401() {
	pwd := "abcdefghijkl"
	user := models.NewUser("test@example.com", "test")
	_ = user.SetPassword(pwd)

	wrongPwd, _ := json.Marshal(&handlers.LoginRequest{Email: user.Email, Password: "wrong"})
	wrongEmail, _ := json.Marshal(&handlers.LoginRequest{Email: "wrong", Password: user.Password})

	testCases := []struct {
		name    string
		payload []byte
	}{
		{"user does not exist", wrongEmail},
		{"wrong password", wrongPwd},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			s.svc.EXPECT().
				FindOneByEmailOrUsername(mock.Anything, mock.Anything, mock.Anything).
				Return(nil, data.ErrNoDocuments)

			s.server.ServeHTTP(resp, req)

			var result echo.HTTPError
			_ = json.Unmarshal(resp.Body.Bytes(), &result)

			assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
			assert.Contains(s.T(), "invalid email or password", result.Message)
		})
	}
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Login_400() {
	b, _ := json.Marshal(&handlers.LoginRequest{})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
}

// TODO: fix
//func (s *AuthHandlerTestSuite) TestAuthHandler_Login_422() {
//	b, _ := json.Marshal(`{"username":"foo","password":"bar","derp":"dep"}`)
//
//	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(b))
//	req.Header.Set("Content-Type", "application/json")
//	resp := httptest.NewRecorder()
//
//	s.server.ServeHTTP(resp, req)
//
//	assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.Code)
//}

func (s *AuthHandlerTestSuite) TestAuthHandler_Logout_204_Cookie() {
	user := models.NewUser("test@example.com", "test")
	_, refresh, _ := user.Login()

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewRefreshTokenCookie(refresh))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil)

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNoContent, resp.Code)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Logout_401_Cookie_Invalid() {
	user := models.NewUser("test@example.com", "test")
	_, _, _ = user.Login()

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewRefreshTokenCookie([]byte("invalid")))
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), jwt.ErrTokenInvalid, result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Logout_401_Cookie_Mismatch() {
	user := models.NewUser("test@example.com", "test")
	_, refresh, _ := user.Login()
	user.Logout()

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewRefreshTokenCookie(refresh))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), "token mismatch", result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Logout_204_Token() {
	user := models.NewUser("test@example.com", "test")
	_, refresh, _ := user.Login()

	payload := &handlers.LogoutRequest{
		RefreshToken: string(refresh),
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil)

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusNoContent, resp.Code)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Logout_400_Body_Missing_Key() {
	b, _ := json.Marshal(&handlers.LoginRequest{})

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.Code)
	assert.Equal(s.T(), jwt.ErrBodyMissingKey, result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Logout_400_Token_Missing() {
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(s.T(), jwt.ErrRequestMalformed, result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Logout_401_Token_Mismatch() {
	user := models.NewUser("test@example.com", "test")
	_, refresh, _ := user.Login()
	user.Logout()

	payload := &handlers.LogoutRequest{
		RefreshToken: string(refresh),
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), "token mismatch", result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Refresh_200_Cookie() {
	user := models.NewUser("test@example.com", "test")
	_, refresh, _ := user.Login()

	payload := &handlers.RefreshRequest{
		GrantType: "refresh_token",
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewRefreshTokenCookie(refresh))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil)

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)

	s.server.ServeHTTP(resp, req)

	var result handlers.RefreshResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	expected := 2
	if viper.GetBool(config.CSRFEnabled) {
		expected = 3
	}
	if assert.Equal(s.T(), expected, len(resp.Result().Cookies())) {
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
		assert.Equal(s.T(), expected, cookies)
	}

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.NotEqual(s.T(), "", result.AccessToken)
	assert.NotEqual(s.T(), "", result.ExpiresIn)
	assert.NotEqual(s.T(), "", result.RefreshToken)
	assert.NotEqual(s.T(), "", result.TokenType)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Refresh_400_Cookie_Missing() {
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(s.T(), jwt.ErrRequestMalformed, result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Refresh_401_Cookie_Invalid() {
	user := models.NewUser("test@example.com", "test")
	_, refresh, _ := user.Login()
	user.Logout()

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewRefreshTokenCookie(refresh))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), "token mismatch", result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Refresh_200_Token() {
	user := models.NewUser("test@example.com", "test")
	_, refresh, _ := user.Login()

	payload := &handlers.RefreshRequest{
		GrantType:    "refresh_token",
		RefreshToken: string(refresh),
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil)

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)

	s.server.ServeHTTP(resp, req)

	var result handlers.RefreshResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	expected := 2
	if viper.GetBool(config.CSRFEnabled) {
		expected = 3
	}
	if assert.Equal(s.T(), expected, len(resp.Result().Cookies())) {
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
		assert.Equal(s.T(), expected, cookies)
	}

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.NotEqual(s.T(), "", result.AccessToken)
	assert.NotEqual(s.T(), "", result.ExpiresIn)
	assert.NotEqual(s.T(), "", result.RefreshToken)
	assert.NotEqual(s.T(), "", result.TokenType)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Refresh_400_Token_Missing() {
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(s.T(), jwt.ErrRequestMalformed, result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Refresh_401_Token_Invalid() {
	payload := &handlers.RefreshRequest{
		GrantType:    "refresh_token",
		RefreshToken: "invalid",
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), "token invalid", result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Refresh_401_Token_Mismatch() {
	user := models.NewUser("test@example.com", "test")
	_, refresh, _ := user.Login()
	user.Logout()

	payload := &handlers.RefreshRequest{
		GrantType:    "refresh_token",
		RefreshToken: string(refresh),
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), "token mismatch", result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Signup_200() {
	payload := &handlers.SignUpRequest{
		Email:    "test@example.com",
		Username: "test",
		Name:     "Test",
		Password: "abcdefghijkl",
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		FindOneByEmailOrUsername(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)

	s.svc.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(&models.User{
			Model:    &models.Model{Id: "1"},
			Email:    payload.Email,
			Username: payload.Username,
			Name:     payload.Name,
		}, nil)

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Signup_409() {
	payload := &handlers.SignUpRequest{
		Email:    "test@example.com",
		Username: "test",
		Name:     "Test",
		Password: "abcdefghijkl",
	}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		FindOneByEmailOrUsername(mock.Anything, mock.Anything, mock.Anything).
		Return(&models.User{
			Model:    &models.Model{Id: "1"},
			Email:    payload.Email,
			Username: payload.Username,
			Name:     payload.Name,
		}, nil)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusConflict, resp.Code)
	assert.Equal(s.T(), "email or username already in-use", result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Signup_422() {
	payload := &handlers.SignUpRequest{}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	assert.Equal(s.T(), http.StatusUnprocessableEntity, resp.Code)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Token_200() {
	user := models.NewUser("test@example.com", "test")
	access, _, _ := user.Login()

	token, _ := util.ParseToken(access)

	req := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result handlers.TokenResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	roles, _ := token.Get("roles")
	typ, _ := token.Get("type")

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), token.Expiration(), result.Exp)
	assert.Equal(s.T(), token.IssuedAt(), result.Iat)
	assert.Equal(s.T(), token.Issuer(), result.Iss)
	assert.Equal(s.T(), token.NotBefore(), result.Nbf)
	assert.ElementsMatch(s.T(), roles, user.Roles)
	assert.Equal(s.T(), token.Subject(), result.Sub)
	assert.Equal(s.T(), typ, result.Type)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Token_401() {
	req := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), "token invalid", result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Cookie_200() {
	user := models.NewUser("test@example.com", "test")
	access, _, _ := user.Login()

	token, _ := util.ParseToken(access)

	req := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewAccessTokenCookie(access))
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result handlers.TokenResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	roles, _ := token.Get("roles")
	typ, _ := token.Get("type")

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), token.Expiration(), result.Exp)
	assert.Equal(s.T(), token.IssuedAt(), result.Iat)
	assert.Equal(s.T(), token.Issuer(), result.Iss)
	assert.Equal(s.T(), token.NotBefore(), result.Nbf)
	assert.ElementsMatch(s.T(), roles, user.Roles)
	assert.Equal(s.T(), token.Subject(), result.Sub)
	assert.Equal(s.T(), typ, result.Type)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Cookie_401() {
	req := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewAccessTokenCookie([]byte("wrong")))
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), "token invalid", result.Message)
}
