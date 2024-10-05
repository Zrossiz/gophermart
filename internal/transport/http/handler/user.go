package handler

import "net/http"

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u *UserHandler) Registration(rw http.ResponseWriter, r *http.Request) {

}

func (u *UserHandler) Login(rw http.ResponseWriter, r *http.Request) {

}
