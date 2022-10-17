package app

import (
	"net/http"

	"github.com/alexferl/golib/http/router"
	"github.com/alexferl/golib/http/server"

	"github.com/alexferl/echo-boilerplate/handlers"
)

func Start() {
	c := NewConfig()
	c.BindFlags()

	s := server.New()
	h := &handlers.Handler{}
	r := &router.Router{
		Routes: []*router.Route{
			{"Root", http.MethodGet, "/", h.Root},
			{"Healthz", http.MethodGet, "/healthz", h.Healthz},
		},
	}

	s.Start(r)
}
