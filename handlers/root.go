package handlers

import (
	"fmt"
	"net/http"

	"github.com/alexferl/echo-openapi"
	"github.com/alexferl/golib/http/api/server"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"github.com/alexferl/echo-boilerplate/config"
)

type RootHandler struct {
	*openapi.Handler
}

func NewRootHandler(openapi *openapi.Handler) *RootHandler {
	return &RootHandler{
		Handler: openapi,
	}
}

func (h *RootHandler) Register(s *server.Server) {
	s.Add(http.MethodGet, "/", h.Root)
}

// Root returns the welcome message.
func (h *RootHandler) Root(c echo.Context) error {
	m := fmt.Sprintf("Welcome to %server", viper.GetString(config.AppName))
	return c.JSON(http.StatusOK, echo.Map{"message": m})
}
