package dto

import "github.com/shopspring/decimal"

type CreateBalanceHistory struct {
	OrderID int64           `json:"order_id"`
	UserID  int64           `json:"user_id"`
	Change  decimal.Decimal `json:"change"`
}

type Registration struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

type CreateStatus struct {
	Name string `json:"name"`
}

type BalanceUser struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withrawn"`
}
