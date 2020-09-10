package handlers_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	test "echo-boilerplate/internal/pkg/testing"
)

func TestHealthz(t *testing.T) {
	e := echo.New()
	_, rec, req := test.NewTestRequest(e, "Healthz", "/healthz", nil)
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}
