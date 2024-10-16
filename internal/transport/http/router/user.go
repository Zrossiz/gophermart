package router

import (
	"net/http"

	"github.com/Zrossiz/gophermart/internal/mIDdleware"
	"github.com/go-chi/chi/v5"
)

type UserRouter struct {
	handler UserHandler
}

type UserHandler interface {
	Login(rw http.ResponseWriter, r *http.Request)
	Registration(rw http.ResponseWriter, r *http.Request)
	UploadOrder(rw http.ResponseWriter, r *http.Request)
	GetAllOrdersByUser(rw http.ResponseWriter, r *http.Request)
	GetUserBalance(rw http.ResponseWriter, r *http.Request)
	Withdraw(rw http.ResponseWriter, r *http.Request)
}

func NewUserRouter(h UserHandler) *UserRouter {
	return &UserRouter{handler: h}
}

func (u *UserRouter) RegisterRoutes(r chi.Router, h UserHandler) {
	r.Route("/api/user", func(r chi.Router) {
		r.Post(("/register"), u.handler.Registration)
		r.Post("/login", u.handler.Login)
		r.With(mIDdleware.JWTMIDdleware).Post("/orders", u.handler.UploadOrder)
		r.With(mIDdleware.JWTMIDdleware).Get("/orders", u.handler.GetAllOrdersByUser)
		r.With(mIDdleware.JWTMIDdleware).Get("/balance", u.handler.GetUserBalance)
		r.With(mIDdleware.JWTMIDdleware).Post("/balance/withdraw", u.handler.Withdraw)
	})
}
