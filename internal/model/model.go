package model

import (
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"password" db:"password"`
	Account   float64   `json:"account" db:"account"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Status struct {
	ID        int       `json:"id" db:"id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Order struct {
	OrderID   string    `json:"number" db:"order_id"`
	UserID    int       `json:"-" db:"user_id"`
	Status    string    `json:"status,omitempty" db:"status"`
	Accrual   float64   `json:"accrual,omitempty" db:"accrual"`
	CreatedAt time.Time `json:"uploaded_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type BalanceHistory struct {
	ID          int        `json:"-" db:"id"`
	OrderID     string     `json:"order" db:"order_id"`
	UserID      int        `json:"-" db:"user_id"`
	Change      float64    `json:"sum" db:"change"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt   time.Time  `json:"-" db:"created_at"`
	UpdatedAt   time.Time  `json:"-" db:"updated_at"`
}

type RefreshToken struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"userid"`
	Token     string    `json:"token" db:"token"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
