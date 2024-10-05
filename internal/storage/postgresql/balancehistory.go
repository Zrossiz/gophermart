package postgresql

import (
	"context"

	"github.com/Zrossiz/gophermart/internal/dto"
	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type BalanceHistoryStore struct {
	db *pgxpool.Pool
}

func NewBalanceHistoryStore(db *pgxpool.Pool) BalanceHistoryStore {
	return BalanceHistoryStore{db: db}
}

func (b *BalanceHistoryStore) Create(balanceHistoryDTO dto.CreateBalanceHistory) (bool, error) {
	sql := `INSERT INTO balance_history (order_id, user_id, change) VALUES ($1, $2, $3)`
	_, err := b.db.Exec(context.Background(), sql, balanceHistoryDTO.OrderID, balanceHistoryDTO.UserID, balanceHistoryDTO.Change)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (b *BalanceHistoryStore) GetAllDebits(userID int64) ([]model.BalanceHistory, error) {
	sql := `SELECT id, order_id, user_id, change, created_id, updated_id FROM balance_history WHERE user_id = $1`
	rows, err := b.db.Query(context.Background(), sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []model.BalanceHistory
	for rows.Next() {
		var history model.BalanceHistory

		err := rows.Scan(&history.ID, &history.OrderID, &history.UserID, &history.Change, &history.CreatedAt, &history.UpdatedAt)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return histories, nil
}
