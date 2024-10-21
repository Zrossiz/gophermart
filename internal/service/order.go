package service

import (
	"fmt"
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
	CreateOrder(orderID string, userID int) (bool, error)
	GetAllOrdersByUser(userID int64) ([]model.Order, error)
	UpdateSumAndStatusOrder(orderID string, status string, sum float64) (bool, error)
	GetOrderByID(orderID string) (*model.Order, error)
	GetAllWithdrawnByUser(userID int64) (float64, error)
	GetAllUnhandlerOrders(unhandledStatus1, unhandledStatus2 int) ([]string, error)
}

type APIService interface {
	UpdateOrder(orderID string) (dto.ExternalOrderResponse, error)
}

func NewOrderService(db OrderStorage, statusDB StatusStorage, a APIService, log *zap.Logger) *OrderService {
	return &OrderService{
		orderDB:  db,
		log:      log,
		api:      a,
		statusDB: statusDB,
	}
}

func (o *OrderService) UploadOrder(order string, userID int) error {
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

	luhn := utils.IsLuhn(order)
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
		o.log.Error("db query error", zap.Error(err))
		return make([]dto.ResponseOrder, 0), apperrors.ErrDBQuery
	}

	if orders == nil {
		return nil, apperrors.ErrOrdersNotFound
	}

	if len(orders) == 0 {
		return nil, apperrors.ErrOrdersNotFound
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

	var unhandledStatuses []int

	for i := 0; i < len(statuses); i++ {
		if statuses[i].Status == "new" || statuses[i].Status == "processing" {
			unhandledStatuses = append(unhandledStatuses, statuses[i].ID)
		}
	}

	unhandledOrders, err := o.orderDB.GetAllUnhandlerOrders(unhandledStatuses[0], unhandledStatuses[1])
	if err != nil {
		o.log.Error("error get unhandledOrders", zap.Error(err))
		return
	}

	if len(unhandledOrders) == 0 {
		return
	}

	for _, order := range unhandledOrders {
		respOrder, err := o.api.UpdateOrder(order)
		if err != nil {
			o.log.Error("error update order from external app", zap.Error(err))
		}

		_, err = o.orderDB.UpdateSumAndStatusOrder(order, strings.ToLower(respOrder.Status), respOrder.Accrual)

		if err != nil {
			o.log.Error("error update order", zap.Error(err))
		}
	}
}
