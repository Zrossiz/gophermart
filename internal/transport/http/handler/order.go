package handler

import "github.com/Zrossiz/gophermart/internal/model"

type OrderHandler struct {
}

type OrderService interface {
	UploadOrder(order int, userId int) error
	GetAllOrdersByUser(userId int) ([]model.Order, error)
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}
