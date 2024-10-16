package handler

type BalanceHistoryHandler struct {
}

type BalanceHistoryService interface {
	Withdraw(userID, orderID, sum int) error
}

func NewBalanceHistoryHandler() *BalanceHistoryHandler {
	return &BalanceHistoryHandler{}
}
