package handlers_test

import (
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	test "echo-boilerplate/internal/pkg/testing"
)

func TestRoot(t *testing.T) {
	e := echo.New()
	_, rec, req := test.NewTestRequest(e, "Root", "/", nil)
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
