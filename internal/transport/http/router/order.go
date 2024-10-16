package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type OrderRouter struct {
	handler OrderHandler
}

type OrderHandler interface {
	UpdateOrders(rw http.ResponseWriter, r *http.Request)
}

func NewOrderRouter(handler OrderHandler) *OrderRouter {
	return &OrderRouter{handler: handler}
}

func (o *OrderRouter) RegisterRoutes(r chi.Router, h OrderHandler) {
	r.Route("/order", func(r chi.Router) {
		r.Post("/update/orders", o.handler.UpdateOrders)
	})
}
