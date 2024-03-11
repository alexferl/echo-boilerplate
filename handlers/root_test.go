package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexferl/echo-openapi"
	"github.com/stretchr/testify/assert"

	"github.com/alexferl/echo-boilerplate/handlers"
	"github.com/alexferl/echo-boilerplate/server"
	_ "github.com/alexferl/echo-boilerplate/testing"
)

func TestHandler_Root(t *testing.T) {
	h := handlers.NewRootHandler(openapi.NewHandler())
	s := server.NewTestServer(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Welcome")
}
