package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Zrossiz/gophermart/internal/dto"
)

type StatusHandler struct {
	service StatusService
}

type StatusService interface {
	Create(name string) error
}

func NewStatusHandler(serv StatusService) *StatusHandler {
	return &StatusHandler{service: serv}
}

func (o *StatusHandler) Create(rw http.ResponseWriter, r *http.Request) {
	var statusDTO dto.CreateStatus

	err := json.NewDecoder(r.Body).Decode(&statusDTO)
	if err != nil {
		http.Error(rw, "invalid request body", http.StatusBadRequest)
		return
	}

	err = o.service.Create(statusDTO.Name)
	if err != nil {
		http.Error(rw, "insert db error", http.StatusInternalServerError)
		return
	}
}
