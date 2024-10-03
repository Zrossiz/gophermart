package dto

import "github.com/shopspring/decimal"

type CreateBalanceHistory struct {
	OrderID int64           `json:"order_id"`
	UserID  int64           `json:"user_id"`
	Change  decimal.Decimal `json:"change"`
}
