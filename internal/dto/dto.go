package dto

import (
	"time"
)

type CreateBalanceHistory struct {
	OrderID string  `json:"order_ID"`
	UserID  int64   `json:"user_ID"`
	Change  float64 `json:"change"`
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
	Withdrawn float64 `json:"withdrawn"`
}

type Withdraw struct {
	Sum   float64 `json:"sum"`
	Order string  `json:"order"`
}

type ResponseOrder struct {
	OrderID   string    `json:"number"`
	Accrual   float64   `json:"accrual,omitempty"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"uploaded_at"`
}

type ExternalOrderResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
