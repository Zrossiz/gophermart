package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserRouter struct{}

type UserHandler interface {
	Login(rw http.ResponseWriter, r *http.Request)
	Registration(rw http.ResponseWriter, r *http.Request)
}

func NewUserRouter() *UserRouter {
	return &UserRouter{}
}

func (u *UserRouter) RegisterRoutes(r chi.Router, h UserHandler) {
	r.Route("/user", func(r chi.Router) {
	})
}
