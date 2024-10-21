package apperrors

import "errors"

var (
	ErrUserAlreadyExists                 = errors.New("user already exists")
	ErrUserNotFound                      = errors.New("user not found")
	ErrInvalIDPassword                   = errors.New("invalID password")
	ErrDBQuery                           = errors.New("database query error")
	ErrHashPassword                      = errors.New("error hashing password")
	ErrJWTGeneration                     = errors.New("error generating JWT")
	ErrSaveToken                         = errors.New("error saving token")
	ErrInvalIDOrderID                    = errors.New("invalID order ID")
	ErrOrderAlreadyUploaded              = errors.New("order already uploaded")
	ErrOrderAlreadyUploadedByAnotherUser = errors.New("order already uploaded by another user")
	ErrOrdersNotFound                    = errors.New("orders not found")
	ErrNotEnoughMoney                    = errors.New("not enough money")
	ErrWithdrawlsNotFound                = errors.New("withdrawsls not found")

	ErrNotFoundOrders      = errors.New("нет заказов для обработки")
	ErrInvalidAccrualOrder = errors.New("invalid accrual order")                      // 500
	ErrNoContent           = errors.New("заказ не зарегистрирован в системе расчёта") // 204
	ErrTooManyRequests     = errors.New("превышено количество запросов к сервису")    // 429
)
