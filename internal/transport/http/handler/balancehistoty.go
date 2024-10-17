package handler

type BalanceHistoryHandler struct {
	service BalanceHistoryService
}

type BalanceHistoryService interface {
	Withdraw(userID, orderID int, sum float64) error
}

func NewBalanceHistoryHandler(serv BalanceHistoryService) *BalanceHistoryHandler {
	return &BalanceHistoryHandler{service: serv}
}
