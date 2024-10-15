package service

import (
	"strconv"

	"github.com/Zrossiz/gophermart/internal/apperrors"
	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/Zrossiz/gophermart/internal/utils"
	"go.uber.org/zap"
)

type OrderService struct {
	db  OrderStorage
	log *zap.Logger
}

type OrderStorage interface {
	CreateOrder(orderId int, userId int) (bool, error)
	GetAllOrdersByUser(userID int64) ([]model.Order, error)
	UpdateStatusOrder(orderID int64, statusID int) (bool, error)
	GetOrderById(orderId int) (*model.Order, error)
	GetAllWithdrawnByUser(userID int64) (float64, error)
}

func NewOrderService(db OrderStorage, log *zap.Logger) *OrderService {
	return &OrderService{
		db:  db,
		log: log,
	}
}

func (o *OrderService) UploadOrder(order int, userId int) error {
	existOrder, err := o.db.GetOrderById(order)
	if err != nil {
		o.log.Error(err.Error())
		return apperrors.ErrDBQuery
	}

	if existOrder != nil && existOrder.UserID != userId {
		return apperrors.ErrOrderAlreadyUploadedByAnotherUser
	}

	if existOrder != nil {
		return apperrors.ErrOrderAlreadyUploaded
	}

	luhn := utils.IsLuhn(strconv.Itoa(order))
	if !luhn {
		return apperrors.ErrInvalidOrderId
	}

	_, err = o.db.CreateOrder(order, userId)
	if err != nil {
		o.log.Error(err.Error())
		return apperrors.ErrDBQuery
	}

	return nil
}

func (o *OrderService) GetAllOrdersByUser(userId int) ([]model.Order, error) {
	orders, err := o.db.GetAllOrdersByUser(int64(userId))
	if err != nil {
		o.log.Error(err.Error())
		return make([]model.Order, 0), apperrors.ErrDBQuery
	}

	if len(orders) == 0 {
		return make([]model.Order, 0), apperrors.ErrOrdersNotFound
	}

	return orders, nil
}
