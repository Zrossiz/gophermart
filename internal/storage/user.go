package storage

import (
	"context"

	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shopspring/decimal"
)

type UserStore struct {
	db *pgxpool.Pool
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

func (u *UserStore) Create(name string, password string) (bool, error) {
	sql := `INSERT INTO "user" (name, password) VALUES ($1, $2)`
	_, err := u.db.Exec(context.Background(), sql, name, password)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserStore) GetUserByName(name string) (*model.User, error) {
	sql := `SELECT id, name, password, account, created_at, updated_at FROM "user" WHERE name = $1`
	row := u.db.QueryRow(context.Background(), sql, name)
	var user model.User
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Account, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserStore) UpdateUserBalance(userID int64, balance decimal.Decimal) (bool, error) {
	sql := `UPDATE "user" SET account = $1 WHERE id = $2`
	cmdTag, err := u.db.Exec(context.Background(), sql, balance, userID)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}
