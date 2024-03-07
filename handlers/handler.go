package handlers

import (
	"github.com/alexferl/golib/http/api/server"
)

type Handler interface {
	Register(s *server.Server)
}
