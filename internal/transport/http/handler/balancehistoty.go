package handler

type BalanceHistoryHandler struct {
}

type BalanceHistoryService interface {
}

func NewBalanceHistoryHandler() *BalanceHistoryHandler {
	return &BalanceHistoryHandler{}
}
