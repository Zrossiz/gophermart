package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DBStorage struct {
	BalanceHistoryStore *BalanceHistoryStore
	UserStore           *UserStore
	OrderStore          *OrderStore
	TokenStore          *TokenStore
	StatusStore         *StatusStore
}

func Connect(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("connect fail")
	}

	return db, nil
}

func New(dbConn *pgxpool.Pool) *DBStorage {

	return &DBStorage{
		BalanceHistoryStore: NewBalanceHistoryStore(dbConn),
		UserStore:           NewUserStore(dbConn),
		OrderStore:          NewOrderStore(dbConn),
		TokenStore:          NewTokenStore(dbConn),
		StatusStore:         NewStatusStore(dbConn),
	}
}
