package service

import "github.com/Zrossiz/gophermart/internal/model"

type OrderService struct {
	db *OrderStorage
}

type OrderStorage interface {
	CreateOrder(order *model.Order) (bool, error)
	GetAllOrdersByUser(userID int64) ([]model.Order, error)
	UpdateStatusOrder(orderID int64, statusID int) (bool, error)
}

func NewOrderService(db *OrderStorage) *OrderService {
	return &OrderService{
		db: db,
	}
}
