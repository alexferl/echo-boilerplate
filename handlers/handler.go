package handlers

import (
	"net/http"

	"github.com/alexferl/golib/http/api/server"
)

type BaseHandler interface {
	AddRoutes(s *server.Server)
}

type Handler struct{}

func NewHandler() BaseHandler {
	return &Handler{}
}

func (h *Handler) AddRoutes(s *server.Server) {
	s.Add(http.MethodGet, "/", h.Root)
}
