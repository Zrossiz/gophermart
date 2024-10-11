package router

import (
	"net/http"

	"github.com/Zrossiz/gophermart/internal/transport/http/handler"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	Order          OrderRouter
	User           UserRouter
	BalanceHistory BalanceHistoryRouter
	Status         StatusRouter
}

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	router := &Router{
		Order:          *NewOrderRouter(),
		User:           *NewUserRouter(h.UserHandler),
		BalanceHistory: *NewBalanceHistoryRouter(),
		Status:         *NewStatusRouter(),
	}

	router.User.RegisterRoutes(r, h.UserHandler)
	router.BalanceHistory.RegisterRoutes(r, h.BalanceHistoryHandler)
	router.Order.RegisterRoutes(r, h.OrderHandler)
	router.Status.RegisterRoutes(r, h.StatusHandler)

	return r
}
