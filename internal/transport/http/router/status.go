package router

import (
	"github.com/go-chi/chi/v5"
)

type StatusRouter struct {
}

type StatusHandler interface {
}

func NewStatusRouter() *StatusRouter {
	return &StatusRouter{}
}

func (s *StatusRouter) RegisterRoutes(r chi.Router, h StatusHandler) {
	r.Route("/status", func(r chi.Router) {
	})
}
