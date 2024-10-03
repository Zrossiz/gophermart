package storage

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type StatusStore struct {
	db *pgxpool.Pool
}

func NewStatusStore(db *pgxpool.Pool) *StatusStore {
	return &StatusStore{db: db}
}
