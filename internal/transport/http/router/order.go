package router

import "github.com/go-chi/chi/v5"

type OrderRouter struct {
}

type OrderHandler interface {
}

func NewOrderRouter() *OrderRouter {
	return &OrderRouter{}
}

func (o *OrderRouter) RegisterRoutes(r chi.Router, h OrderHandler) {
	r.Route("/order", func(r chi.Router) {
	})
}
