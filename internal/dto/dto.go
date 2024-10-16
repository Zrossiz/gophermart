package dto

import "github.com/shopspring/decimal"

type CreateBalanceHistory struct {
	OrderID int64           `json:"order_ID"`
	UserID  int64           `json:"user_ID"`
	Change  decimal.Decimal `json:"change"`
}

type Registration struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CreateStatus struct {
	Name string `json:"name"`
}

type BalanceUser struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withrawn"`
}

type Withdraw struct {
	Sum   float64 `json:"sum"`
	Order float64 `json:"order"`
}
