package handler

type BalanceHistoryHandler struct {
}

type BalanceHistoryService interface {
	Withdraw(userId, orderId, sum int) error
}

func NewBalanceHistoryHandler() *BalanceHistoryHandler {
	return &BalanceHistoryHandler{}
}
