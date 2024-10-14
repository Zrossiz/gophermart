package service

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/Zrossiz/gophermart/internal/apperrors"
	"github.com/Zrossiz/gophermart/internal/model"
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
}

func NewOrderService(db OrderStorage, log *zap.Logger) *OrderService {
	return &OrderService{
		db:  db,
		log: log,
	}
}

func (o *OrderService) UploadOrder(order int, userId int) error {
	luhn := isLuhn(strconv.Itoa(order))
	fmt.Println("luhn: ", luhn)
	if !luhn {
		return apperrors.ErrInvalidOrderId
	}

	_, err := o.db.CreateOrder(order, userId)
	if err != nil {
		o.log.Error(err.Error())
		return apperrors.ErrDBQuery
	}

	return nil
}

func isLuhn(orderId string) bool {
	var sum int
	var alternate bool

	for i := len(orderId) - 1; i >= 0; i-- {
		r := rune(orderId[i])

		if !unicode.IsDigit(r) {
			return false
		}

		n, _ := strconv.Atoi(string(r))

		if alternate {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}

		sum += n
		alternate = !alternate
	}

	return sum%10 == 0
}
