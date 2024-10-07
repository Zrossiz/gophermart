package handler

type OrderHandler struct {
}

type OrderService interface {
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}
