package router

import (
	"net/http"

	"github.com/Zrossiz/gophermart/internal/middleware"
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

func (s *StatusRouter) RegisterRoutes(r chi.Router, _ StatusHandler) {
	r.Route("/api/statuses", func(r chi.Router) {
		r.With(middleware.JWTMiddleware).Post("/", s.handler.Create)
	})
}
