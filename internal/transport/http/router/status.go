package router

import (
	"net/http"

	"github.com/Zrossiz/gophermart/internal/mIDdleware"
	"github.com/go-chi/chi/v5"
)

type StatusRouter struct {
	handler StatusHandler
}

type StatusHandler interface {
	Create(rw http.ResponseWriter, r *http.Request)
}

func NewStatusRouter(handler StatusHandler) *StatusRouter {
	return &StatusRouter{
		handler: handler,
	}
}

func (s *StatusRouter) RegisterRoutes(r chi.Router, h StatusHandler) {
	r.Route("/api/statuses", func(r chi.Router) {
		r.With(mIDdleware.JWTMIDdleware).Post("/", s.handler.Create)
	})
}
