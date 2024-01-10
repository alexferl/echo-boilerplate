package handlers

import (
	"net/http"

	"github.com/alexferl/golib/http/api/server"
)

type IHandler interface {
	AddRoutes(s *server.Server)
}

type Handler struct{}

func NewHandler() IHandler {
	return &Handler{}
}

func (h *Handler) AddRoutes(s *server.Server) {
	s.Add(http.MethodGet, "/", h.Root)
}
