package handler

import (
	"net/http"

	"github.com/Zrossiz/gophermart/internal/dto"
)

type OrderHandler struct {
	service OrderService
}

type OrderService interface {
	UploadOrder(order int, userID int) error
	GetAllOrdersByUser(userID int) ([]dto.ResponseOrder, error)
	UpdateOrders()
}

func NewOrderHandler(serv OrderService) *OrderHandler {
	return &OrderHandler{service: serv}
}

func (o *OrderHandler) UpdateOrders(rw http.ResponseWriter, r *http.Request) {
	o.service.UpdateOrders()
	// if err != nil {
	// 	http.Error(rw, "unknown error", http.StatusInternalServerError)
	// 	return
	// }

	rw.WriteHeader(http.StatusOK)
}
