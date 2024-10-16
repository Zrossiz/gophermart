package postgresql

import (
	"context"

	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type StatusStore struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func NewStatusStore(db *pgxpool.Pool, log *zap.Logger) StatusStore {
	return StatusStore{db: db, log: log}
}

func (s *StatusStore) Create(status string) (bool, error) {
	sql := `INSERT INTO statuses (status) VALUES ($1)`
	_, err := s.db.Exec(context.Background(), sql, status)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *StatusStore) GetAll() ([]model.Status, error) {
	sql := `SELECT ID, status, created_at, updated_at FROM statuses`
	rows, err := s.db.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []model.Status
	for rows.Next() {
		var status model.Status
		err := rows.Scan(&status.ID, &status.Status, &status.CreatedAt, &status.UpdatedAt)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return statuses, nil
}
