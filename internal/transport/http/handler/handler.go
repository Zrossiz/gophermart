package handler

type Handler struct {
	UserHandler           *UserHandler
	BalanceHistoryHandler *BalanceHistoryHandler
	OrderHandler          *OrderHandler
	StatusHandler         *StatusHanlder
}

type Service struct {
}

func New() *Handler {
	return &Handler{
		UserHandler:           NewUserHandler(),
		BalanceHistoryHandler: NewBalanceHistoryHandler(),
		StatusHandler:         NewStatusHandler(),
		OrderHandler:          NewOrderHandler(),
	}
}
