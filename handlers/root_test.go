package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	app "github.com/alexferl/echo-boilerplate"
	"github.com/alexferl/echo-boilerplate/config"
)

var overrides = map[string]any{
	config.JWTPrivateKey:   "../private-key.pem",
	config.CasbinModel:     "../casbin/model.conf",
	config.CasbinPolicy:    "../casbin/policy.csv",
	config.OpenAPISchema:   "../openapi/openapi.yaml",
	config.HTTPLogRequests: false,
}

func TestHandler_Root(t *testing.T) {
	s := app.NewServerWithOverrides(overrides, app.DefaultHandlers()...)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Welcome")
}
