package storage

import (
	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserStore struct {
	db *pgxpool.Pool
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

func (u *UserStore) Create() (bool, error) {
	return true, nil
}

func (u *UserStore) GetUserByName(username string) (*model.User, error) {
	return nil, nil
}
