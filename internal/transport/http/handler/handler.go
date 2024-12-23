package handler

import "go.uber.org/zap"

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

func New(serv Service, log *zap.Logger) *Handler {
	return &Handler{
		UserHandler:           NewUserHandler(serv.UserService, serv.OrderService, serv.BalanceHistoryService, log),
		BalanceHistoryHandler: NewBalanceHistoryHandler(serv.BalanceHistoryService),
		StatusHandler:         NewStatusHandler(serv.StatusService),
		OrderHandler:          NewOrderHandler(serv.OrderService),
	}
}
