package storage

import (
	"context"

	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TokenStore struct {
	db *pgxpool.Pool
}

func NewTokenStore(db *pgxpool.Pool) *TokenStore {
	return &TokenStore{db: db}
}

func (t *TokenStore) Create(userID int64, token string) (bool, error) {
	sql := `INSERT INTO refresh_token (user_id, token) VALUES ($1, $2)`
	_, err := t.db.Exec(context.Background(), sql, userID, token)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t *TokenStore) DeleteByToken(token string) (bool, error) {
	sql := `DELETE FROM refresh_token WHERE token = $1`
	cmdTag, err := t.db.Exec(context.Background(), sql, token)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}

	return true, nil
}

func (t *TokenStore) DeleteTokensByUser(userID int64) (bool, error) {
	sql := `DELETE FROM refresh_token WHERE user_id = $1`
	cmdTag, err := t.db.Exec(context.Background(), sql, userID)
	if err != nil {
		return false, err
	}

	if cmdTag.RowsAffected() == 0 {
		return false, err
	}

	return true, nil
}

func (t *TokenStore) GetTokenByToken(token string) (*model.RefreshToken, error) {
	sql := `SELECT * FROM refresh_token WHERE token = $1`
	row := t.db.QueryRow(context.Background(), sql, token)
	var rt model.RefreshToken
	err := row.Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.CreatedAt, &rt.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &rt, nil
}

func (t *TokenStore) GetTokenByUser(userID int64) (*model.RefreshToken, error) {
	sql := `SELECT * FROM refresh_token WHERE user_id = $1`
	row := t.db.QueryRow(context.Background(), sql, userID)
	var rt model.RefreshToken
	err := row.Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.CreatedAt, &rt.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &rt, nil
}
