package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Zrossiz/gophermart/internal/apperrors"
	"github.com/Zrossiz/gophermart/internal/dto"
	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/Zrossiz/gophermart/internal/utils"
	"go.uber.org/zap"
)

type OrderService struct {
	orderDB  OrderStorage
	statusDB StatusStorage
	log      *zap.Logger
	api      APIService
}

type OrderStorage interface {
	CreateOrder(orderID int, userID int) (bool, error)
	GetAllOrdersByUser(userID int64) ([]model.Order, error)
	UpdateSumAndStatusOrder(orderID int64, status string, sum float64) (bool, error)
	GetOrderByID(orderID int) (*model.Order, error)
	GetAllWithdrawnByUser(userID int64) (float64, error)
	GetAllUnhandlerOrders(unhandledStatus1, unhandledStatus2 int) ([]model.Order, error)
}

type APIService interface {
	UpdateOrder(orderID int) (string, float64, error)
}

func NewOrderService(db OrderStorage, statusDB StatusStorage, a APIService, log *zap.Logger) *OrderService {
	return &OrderService{
		orderDB:  db,
		log:      log,
		api:      a,
		statusDB: statusDB,
	}
}

func (o *OrderService) UploadOrder(order int, userID int) error {
	existOrder, err := o.orderDB.GetOrderByID(order)
	if err != nil {
		o.log.Error(err.Error())
		return apperrors.ErrDBQuery
	}

	if existOrder != nil && existOrder.UserID != userID {
		return apperrors.ErrOrderAlreadyUploadedByAnotherUser
	}

	if existOrder != nil {
		return apperrors.ErrOrderAlreadyUploaded
	}

	luhn := utils.IsLuhn(strconv.Itoa(order))
	if !luhn {
		return apperrors.ErrInvalIDOrderID
	}

	_, err = o.orderDB.CreateOrder(order, userID)
	if err != nil {
		o.log.Error(err.Error())
		return apperrors.ErrDBQuery
	}

	return nil
}

func (o *OrderService) GetAllOrdersByUser(userID int) ([]dto.ResponseOrder, error) {
	orders, err := o.orderDB.GetAllOrdersByUser(int64(userID))
	if err != nil {
		o.log.Error(err.Error())
		return make([]dto.ResponseOrder, 0), apperrors.ErrDBQuery
	}

	if len(orders) == 0 {
		return make([]dto.ResponseOrder, 0), apperrors.ErrOrdersNotFound
	}

	var responseOrders []dto.ResponseOrder

	for _, order := range orders {
		responseOrders = append(responseOrders, dto.ResponseOrder{
			OrderID:   fmt.Sprint(order.OrderID),
			Status:    strings.ToUpper(order.Status),
			Accrual:   order.Accrual,
			CreatedAt: order.CreatedAt,
		})
	}

	return responseOrders, nil
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
