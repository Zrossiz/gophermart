package router

import (
	"net/http"

	"github.com/Zrossiz/gophermart/internal/transport/http/handler"
	"github.com/go-chi/chi/v5"
)

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	return r
}
