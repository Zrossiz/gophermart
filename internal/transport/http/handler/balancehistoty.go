package handler

import "github.com/Zrossiz/gophermart/internal/model"

type BalanceHistoryHandler struct {
	service BalanceHistoryService
}

type BalanceHistoryService interface {
	Withdraw(userID int, orderID string, sum float64) error
	GetAllWithdrawlsByUser(userID int) ([]model.BalanceHistory, error)
}

func NewBalanceHistoryHandler(serv BalanceHistoryService) *BalanceHistoryHandler {
	return &BalanceHistoryHandler{service: serv}
}
