package service

import (
	"strconv"
	"strings"

	"github.com/Zrossiz/gophermart/internal/apperrors"
	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/Zrossiz/gophermart/internal/utils"
	"go.uber.org/zap"
)

type OrderService struct {
	orderDB  OrderStorage
	statusDB StatusStorage
	log      *zap.Logger
	api      ApiService
}

type OrderStorage interface {
	CreateOrder(orderId int, userId int) (bool, error)
	GetAllOrdersByUser(userID int64) ([]model.Order, error)
	UpdateSumAndStatusOrder(orderID int64, status string, sum float64) (bool, error)
	GetOrderById(orderId int) (*model.Order, error)
	GetAllWithdrawnByUser(userID int64) (float64, error)
	GetAllUnhandlerOrders(unhandledStatus1, unhandledStatus2 int) ([]model.Order, error)
}

type ApiService interface {
	UpdateOrder(orderId int) (string, float64, error)
}

func NewOrderService(db OrderStorage, statusDB StatusStorage, a ApiService, log *zap.Logger) *OrderService {
	return &OrderService{
		orderDB:  db,
		log:      log,
		api:      a,
		statusDB: statusDB,
	}
}

func (o *OrderService) UploadOrder(order int, userId int) error {
	existOrder, err := o.orderDB.GetOrderById(order)
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

	_, err = o.orderDB.CreateOrder(order, userId)
	if err != nil {
		o.log.Error(err.Error())
		return apperrors.ErrDBQuery
	}

	return nil
}

func (o *OrderService) GetAllOrdersByUser(userId int) ([]model.Order, error) {
	orders, err := o.orderDB.GetAllOrdersByUser(int64(userId))
	if err != nil {
		o.log.Error(err.Error())
		return make([]model.Order, 0), apperrors.ErrDBQuery
	}

	if len(orders) == 0 {
		return make([]model.Order, 0), apperrors.ErrOrdersNotFound
	}

	return orders, nil
}

func (o *OrderService) UpdateOrders() {
	statuses, err := o.statusDB.GetAll()
	if err != nil {
		o.log.Error("error get all statuses", zap.Error(err))
	}

	unhandledStatuses := make([]int, 2)

	for i := 0; i < len(statuses); i++ {
		if statuses[i].Status == "new" || statuses[i].Status == "processing" {
			unhandledStatuses = append(unhandledStatuses, statuses[i].ID)
		}
	}

	unhandledOrders, err := o.orderDB.GetAllUnhandlerOrders(unhandledStatuses[0], unhandledStatuses[1])
	if err != nil {
		o.log.Error("error get unhandledOrders", zap.Error(err))
	}

	for _, order := range unhandledOrders {
		status, accrual, err := o.api.UpdateOrder(order.OrderID)
		if err != nil {
			o.log.Error("error update order from external app", zap.Error(err))
		}

		_, err = o.orderDB.UpdateSumAndStatusOrder(int64(order.OrderID), strings.ToLower(status), accrual)
		if err != nil {
			o.log.Error("error update order", zap.Error(err))
		}
	}
}
