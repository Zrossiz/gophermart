package storage

import (
	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BalanceHistoryStore struct {
	db *pgxpool.Pool
}

func NewBalanceHistoryStore(db *pgxpool.Pool) *BalanceHistoryStore {
	return &BalanceHistoryStore{db: db}
}

func (b *BalanceHistoryStore) Create() (bool, error) {
	return true, nil
}

func (b *BalanceHistoryStore) GetUserByName(username string) (*model.User, error) {
	return nil, nil
}
