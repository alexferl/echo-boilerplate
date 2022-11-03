package handlers

import (
	"net/http"

	"github.com/alexferl/golib/http/handler"
	"github.com/alexferl/golib/http/router"
)

// Handler represents the structure of our resource
type Handler struct{}

func NewHandler() handler.Handler {
	return &Handler{}
}

func (h *Handler) GetRoutes() []*router.Route {
	return []*router.Route{
		{Name: "Root", Method: http.MethodGet, Pattern: "/", HandlerFunc: h.Root},
		{Name: "Healthz", Method: http.MethodGet, Pattern: "/healthz", HandlerFunc: h.Healthz},
	}
}
