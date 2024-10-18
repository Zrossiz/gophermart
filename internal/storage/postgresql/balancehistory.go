package postgresql

import (
	"context"
	"fmt"
	"time"

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
	sql := `INSERT INTO balance_history (order_id, user_id, change, processed_at) VALUES ($1, $2, $3, $4)`
	_, err := b.db.Exec(
		context.Background(),
		sql,
		balanceHistoryDTO.OrderID,
		balanceHistoryDTO.UserID,
		balanceHistoryDTO.Change,
		time.Now(),
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (b *BalanceHistoryStore) GetAllDebits(userID int64) ([]model.BalanceHistory, error) {
	sql := `SELECT id, order_id, user_id, change, processed_at, created_at, updated_at FROM balance_history WHERE user_id = $1`
	rows, err := b.db.Query(context.Background(), sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []model.BalanceHistory
	for rows.Next() {
		var history model.BalanceHistory

		err := rows.Scan(&history.ID, &history.OrderID, &history.UserID, &history.Change, &history.ProcessedAt, &history.CreatedAt, &history.UpdatedAt)
		if err != nil {
			return nil, err
		}

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			return nil, fmt.Errorf("ошибка загрузки временной зоны: %w", err)
		}

		if history.ProcessedAt != nil {
			processedAt := history.ProcessedAt.In(loc)
			history.ProcessedAt = &processedAt
		}

		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(histories) == 0 {
		return nil, apperrors.ErrWithdrawlsNotFound
	}

	return histories, nil
}

func (b *BalanceHistoryStore) Withdraw(userID, orderID int, sum float64) error {
	tx, err := b.db.Begin(context.Background())
	if err != nil {
		b.log.Error("failed to start transaction")
		return err
	}
	defer tx.Rollback(context.Background())

	var currentBalance float64
	checkBalanceSQL := `SELECT account FROM users WHERE ID = $1`
	err = tx.QueryRow(context.Background(), checkBalanceSQL, userID).Scan(&currentBalance)
	if err != nil {
		b.log.Error("failed to check balance", zap.Error(err))
		return err
	}

	resAmount := currentBalance - sum

	if resAmount < 0 {
		return apperrors.ErrNotEnoughMoney
	}

	updateBalanceSQL := `UPDATE users SET account = $1 WHERE id = $2`
	_, err = tx.Exec(context.Background(), updateBalanceSQL, resAmount, userID)
	if err != nil {
		b.log.Error("failed to insert balance history", zap.Error(err))
		return err
	}

	_, err = b.Create(dto.CreateBalanceHistory{
		OrderID: int64(orderID),
		UserID:  int64(userID),
		Change:  sum,
	})

	if err != nil {
		b.log.Error("failed to create balance history", zap.Error(err))
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		b.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}
