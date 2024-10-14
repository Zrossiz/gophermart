package handler

type Handler struct {
	UserHandler           *UserHandler
	BalanceHistoryHandler *BalanceHistoryHandler
	OrderHandler          *OrderHandler
	StatusHandler         *StatusHandler
}

type Service struct {
	UserService           UserService
	OrderService          OrderService
	StatusService         StatusService
	BalanceHistoryService BalanceHistoryService
}

func New(serv Service) *Handler {
	return &Handler{
		UserHandler:           NewUserHandler(serv.UserService, serv.OrderService),
		BalanceHistoryHandler: NewBalanceHistoryHandler(),
		StatusHandler:         NewStatusHandler(serv.StatusService),
		OrderHandler:          NewOrderHandler(),
	}
}
