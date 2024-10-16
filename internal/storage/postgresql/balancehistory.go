package postgresql

import (
	"context"

	"github.com/Zrossiz/gophermart/internal/apperrors"
	"github.com/Zrossiz/gophermart/internal/dto"
	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type BalanceHistoryStore struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func NewBalanceHistoryStore(db *pgxpool.Pool, log *zap.Logger) BalanceHistoryStore {
	return BalanceHistoryStore{db: db, log: log}
}

func (b *BalanceHistoryStore) Create(balanceHistoryDTO dto.CreateBalanceHistory) (bool, error) {
	sql := `INSERT INTO balance_history (order_ID, user_ID, change) VALUES ($1, $2, $3)`
	_, err := b.db.Exec(context.Background(), sql, balanceHistoryDTO.OrderID, balanceHistoryDTO.UserID, balanceHistoryDTO.Change)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (b *BalanceHistoryStore) GetAllDebits(userID int64) ([]model.BalanceHistory, error) {
	sql := `SELECT ID, order_ID, user_ID, change, created_ID, updated_ID FROM balance_history WHERE user_ID = $1`
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

func (b *BalanceHistoryStore) Withdraw(userID, orderID, sum int) error {
	tx, err := b.db.Begin(context.Background())
	if err != nil {
		b.log.Error("failed to start transaction")
		return err
	}
	defer tx.Rollback(context.Background())

	var currentBalance int
	checkBalanceSQL := `SELECT account FROM users WHERE ID = $1`
	err = tx.QueryRow(context.Background(), checkBalanceSQL, userID).Scan(&currentBalance)
	if err != nil {
		b.log.Error("failed to check baalnce", zap.Error(err))
		return err
	}

	resIDualAmount := currentBalance - sum

	if resIDualAmount < 0 {
		return apperrors.ErrNotEnoughMoney
	}

	updateBalanceSQL := `UPDATE users SET balance = $1 WHERE ID = $2`
	_, err = tx.Exec(context.Background(), updateBalanceSQL, sum, userID)
	if err != nil {
		b.log.Error("failed to insert balance history", zap.Error(err))
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		b.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}
