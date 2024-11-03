package router

import "github.com/go-chi/chi/v5"

type BalanceHistoryRouter struct {
}

type BalanceHistoryHandler interface {
}

func NewBalanceHistoryRouter() *BalanceHistoryRouter {
	return &BalanceHistoryRouter{}
}

func (o *BalanceHistoryRouter) RegisterRoutes(r chi.Router, _ BalanceHistoryHandler) {
	r.Route("/balance-history", func(r chi.Router) {
	})
}
