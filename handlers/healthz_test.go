package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	app "github.com/alexferl/echo-boilerplate"
)

func TestHandler_Healthz(t *testing.T) {
	s := app.NewServerWithOverrides(overrides, app.DefaultHandlers()...)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, resp.Body.String(), "ok")
}
