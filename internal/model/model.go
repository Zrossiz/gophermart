package model

import (
	"time"

	"google.golang.org/genproto/googleapis/type/decimal"
)

type User struct {
	ID        int             `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	Password  string          `json:"password" db:"password"`
	Account   decimal.Decimal `json:"account" db:"account"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type Status struct {
	ID        int       `json:"id" db:"id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Order struct {
	OrderID     int             `json:"order_id" db:"order_id"`
	UserID      int             `json:"user_id" db:"user_id"`
	StatusID    int             `json:"status_id" db:"status_id"`
	Accrual     decimal.Decimal `json:"accrual" db:"accrual"`
	ProcessedAt time.Time       `json:"processed_at" db:"processed_at"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type BalanceHistory struct {
	ID        int             `json:"id" db:"id"`
	OrderID   int             `json:"order_id" db:"order_id"`
	UserID    int             `json:"user_id" db:"user_id"`
	Change    decimal.Decimal `json:"change" db:"change"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type RefreshToken struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
