package router

import (
	"net/http"

	"github.com/Zrossiz/gophermart/internal/middleware"
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
	Withdrawls(rw http.ResponseWriter, r *http.Request)
}

func NewUserRouter(h UserHandler) *UserRouter {
	return &UserRouter{handler: h}
}

func (u *UserRouter) RegisterRoutes(r chi.Router, h UserHandler) {
	r.Route("/api/user", func(r chi.Router) {
		r.Post(("/register"), u.handler.Registration)
		r.Post("/login", u.handler.Login)
		r.With(middleware.JWTMiddleware).Post("/orders", u.handler.UploadOrder)
		r.With(middleware.JWTMiddleware).Get("/orders", u.handler.GetAllOrdersByUser)
		r.With(middleware.JWTMiddleware).Get("/balance", u.handler.GetUserBalance)
		r.With(middleware.JWTMiddleware).Post("/balance/withdraw", u.handler.Withdraw)
		r.With(middleware.JWTMiddleware).Get("/withdrawls", u.handler.Withdrawls)
	})
}
