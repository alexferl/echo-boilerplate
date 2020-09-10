package handlers_test

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"echo-boilerplate/internal/app/httpd/handlers"
)

func TestRegister(t *testing.T) {
	e := echo.New()
	handlers.Register(e)

	assert.Equal(t, len(handlers.Router.Routes), len(e.Routes()))
}
