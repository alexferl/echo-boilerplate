package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtMw "github.com/alexferl/echo-jwt"
	"github.com/alexferl/echo-openapi"
	api "github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/services"
	"github.com/alexferl/echo-boilerplate/util/cookie"
	"github.com/alexferl/echo-boilerplate/util/jwt"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	svc    *handlers.MockUserService
	server *api.Server
}

func (s *AuthHandlerTestSuite) SetupTest() {
	svc := handlers.NewMockUserService(s.T())
	patSvc := handlers.NewMockPersonalAccessTokenService(s.T())
	h := handlers.NewAuthHandler(openapi.NewHandler(), svc)
	s.svc = svc
	s.server = getServer(svc, patSvc, h)
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
			if s.Assert().Equal(expected, len(resp.Result().Cookies())) {
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
				s.Assert().Equal(expected, cookies)
			}

			s.Assert().Equal(http.StatusOK, resp.Code)
			s.Assert().NotEqual("", result.AccessToken)
			s.Assert().NotEqual("", result.ExpiresIn)
			s.Assert().NotEqual("", result.RefreshToken)
			s.Assert().NotEqual("", result.TokenType)
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
				Return(nil, &services.Error{
					Kind: services.NotExist,
				})

			s.server.ServeHTTP(resp, req)

			var result echo.HTTPError
			_ = json.Unmarshal(resp.Body.Bytes(), &result)

			s.Assert().Equal(http.StatusUnauthorized, resp.Code)
			s.Assert().Contains("invalid email or password", result.Message)
		})
	}
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Login_400() {
	b, _ := json.Marshal(&handlers.LoginRequest{})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusBadRequest, resp.Code)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Logout_204_Cookie() {
	user := models.NewUser("test@example.com", "test")
	_, refresh, _ := user.Login()

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie.NewRefreshToken(refresh))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil)

	s.svc.EXPECT().
		Update(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusNoContent, resp.Code)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Logout_401_Cookie_Invalid() {
	user := models.NewUser("test@example.com", "test")
	_, _, _ = user.Login()

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie.NewRefreshToken([]byte("invalid")))
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
	s.Assert().Equal(jwtMw.ErrTokenInvalid, result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Logout_401_Cookie_Mismatch() {
	user := models.NewUser("test@example.com", "test")
	_, refresh, _ := user.Login()
	user.Logout()

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie.NewRefreshToken(refresh))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
	s.Assert().Equal("token mismatch", result.Message)
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

	s.Assert().Equal(http.StatusNoContent, resp.Code)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Logout_400_Body_Missing_Key() {
	b, _ := json.Marshal(&handlers.LoginRequest{})

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusUnprocessableEntity, resp.Code)
	s.Assert().Equal(jwtMw.ErrBodyMissingKey, result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Logout_400_Token_Missing() {
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusBadRequest, resp.Code)
	s.Assert().Equal(jwtMw.ErrRequestMalformed, result.Message)
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

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
	s.Assert().Equal("token mismatch", result.Message)
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
	req.AddCookie(cookie.NewRefreshToken(refresh))
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
	if s.Assert().Equal(expected, len(resp.Result().Cookies())) {
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
		s.Assert().Equal(expected, cookies)
	}

	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().NotEqual("", result.AccessToken)
	s.Assert().NotEqual("", result.ExpiresIn)
	s.Assert().NotEqual("", result.RefreshToken)
	s.Assert().NotEqual("", result.TokenType)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Refresh_400_Cookie_Missing() {
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusBadRequest, resp.Code)
	s.Assert().Equal(jwtMw.ErrRequestMalformed, result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Refresh_401_Cookie_Invalid() {
	user := models.NewUser("test@example.com", "test")
	_, refresh, _ := user.Login()
	user.Logout()

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie.NewRefreshToken(refresh))
	resp := httptest.NewRecorder()

	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil)

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
	s.Assert().Equal("token mismatch", result.Message)
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
	if s.Assert().Equal(expected, len(resp.Result().Cookies())) {
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
		s.Assert().Equal(expected, cookies)
	}

	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().NotEqual("", result.AccessToken)
	s.Assert().NotEqual("", result.ExpiresIn)
	s.Assert().NotEqual("", result.RefreshToken)
	s.Assert().NotEqual("", result.TokenType)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Refresh_400_Token_Missing() {
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusBadRequest, resp.Code)
	s.Assert().Equal(jwtMw.ErrRequestMalformed, result.Message)
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

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
	s.Assert().Equal("token invalid", result.Message)
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

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
	s.Assert().Equal("token mismatch", result.Message)
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

	s.Assert().Equal(http.StatusOK, resp.Code)
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
		Return(nil, &services.Error{
			Kind:    services.Exist,
			Message: services.ErrUserExist.Error(),
		})

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusConflict, resp.Code)
	s.Assert().Equal(services.ErrUserExist.Error(), result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Signup_422() {
	payload := &handlers.SignUpRequest{}
	b, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	s.Assert().Equal(http.StatusUnprocessableEntity, resp.Code)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Token_200() {
	user := models.NewUser("test@example.com", "test")
	access, _, _ := user.Login()
	token, _ := jwt.ParseEncoded(access)

	req := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	// middleware
	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result handlers.TokenResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	typ, _ := token.Get("type")

	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().Equal(token.Expiration(), result.Exp)
	s.Assert().Equal(token.IssuedAt(), result.Iat)
	s.Assert().Equal(token.Issuer(), result.Iss)
	s.Assert().Equal(token.NotBefore(), result.Nbf)
	s.Assert().Equal(token.Subject(), result.Sub)
	s.Assert().Equal(typ, result.Type)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Token_401() {
	req := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
	s.Assert().Equal("token invalid", result.Message)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Cookie_200() {
	user := models.NewUser("test@example.com", "test")
	access, _, _ := user.Login()
	token, _ := jwt.ParseEncoded(access)

	req := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie.NewAccessToken(access))
	resp := httptest.NewRecorder()

	// middleware
	s.svc.EXPECT().
		Read(mock.Anything, mock.Anything).
		Return(user, nil).Once()

	s.server.ServeHTTP(resp, req)

	var result handlers.TokenResponse
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	typ, _ := token.Get("type")

	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().Equal(token.Expiration(), result.Exp)
	s.Assert().Equal(token.IssuedAt(), result.Iat)
	s.Assert().Equal(token.Issuer(), result.Iss)
	s.Assert().Equal(token.NotBefore(), result.Nbf)
	s.Assert().Equal(token.Subject(), result.Sub)
	s.Assert().Equal(typ, result.Type)
}

func (s *AuthHandlerTestSuite) TestAuthHandler_Cookie_401() {
	req := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie.NewAccessToken([]byte("wrong")))
	resp := httptest.NewRecorder()

	s.server.ServeHTTP(resp, req)

	var result echo.HTTPError
	_ = json.Unmarshal(resp.Body.Bytes(), &result)

	s.Assert().Equal(http.StatusUnauthorized, resp.Code)
	s.Assert().Equal("token invalid", result.Message)
}
