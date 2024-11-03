package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type DBStorage struct {
	BalanceHistoryStore BalanceHistoryStore
	UserStore           UserStore
	OrderStore          OrderStore
	TokenStore          TokenStore
	StatusStore         StatusStore
}

func Connect(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("connect fail")
	}

	return db, nil
}

func New(dbConn *pgxpool.Pool, log *zap.Logger) DBStorage {

	return DBStorage{
		BalanceHistoryStore: NewBalanceHistoryStore(dbConn, log),
		UserStore:           NewUserStore(dbConn, log),
		OrderStore:          NewOrderStore(dbConn, log),
		TokenStore:          NewTokenStore(dbConn, log),
		StatusStore:         NewStatusStore(dbConn, log),
	}
}
