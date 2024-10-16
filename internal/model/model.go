package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID        int             `json:"ID" db:"ID"`
	Name      string          `json:"name" db:"name"`
	Password  string          `json:"password" db:"password"`
	Account   decimal.Decimal `json:"account" db:"account"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type Status struct {
	ID        int       `json:"ID" db:"ID"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Order struct {
	OrderID     int             `json:"order_ID" db:"order_ID"`
	UserID      int             `json:"user_ID" db:"user_ID"`
	Status      string          `json:"status,omitempty" db:"status"`
	Accrual     decimal.Decimal `json:"accrual" db:"accrual"`
	ProcessedAt *time.Time      `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type BalanceHistory struct {
	ID        int             `json:"ID" db:"ID"`
	OrderID   int             `json:"order_ID" db:"order_ID"`
	UserID    int             `json:"user_ID" db:"user_ID"`
	Change    decimal.Decimal `json:"change" db:"change"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type RefreshToken struct {
	ID        int       `json:"ID" db:"ID"`
	UserID    int       `json:"user_ID" db:"user_ID"`
	Token     string    `json:"token" db:"token"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
