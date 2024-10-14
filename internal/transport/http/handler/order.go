package handler

type OrderHandler struct {
}

type OrderService interface {
	UploadOrder(order int, userId int) error
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}
