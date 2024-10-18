package model

import (
	"time"
)

type User struct {
	ID        int       `json:"ID" db:"ID"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"password" db:"password"`
	Account   float64   `json:"account" db:"account"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Status struct {
	ID        int       `json:"ID" db:"ID"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Order struct {
	OrderID   int       `json:"number" db:"order_ID"`
	UserID    int       `json:"-" db:"user_ID"`
	Status    string    `json:"status,omitempty" db:"status"`
	Accrual   float64   `json:"accrual,omitempty" db:"accrual"`
	CreatedAt time.Time `json:"uploaded_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type BalanceHistory struct {
	ID          int        `json:"-" db:"ID"`
	OrderID     int        `json:"order" db:"order_ID"`
	UserID      int        `json:"-" db:"user_ID"`
	Change      float64    `json:"sum" db:"change"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	CreatedAt   time.Time  `json:"-" db:"created_at"`
	UpdatedAt   time.Time  `json:"-" db:"updated_at"`
}

type RefreshToken struct {
	ID        int       `json:"ID" db:"ID"`
	UserID    int       `json:"user_ID" db:"user_ID"`
	Token     string    `json:"token" db:"token"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
