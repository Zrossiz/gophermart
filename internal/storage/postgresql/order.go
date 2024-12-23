package postgresql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/Zrossiz/gophermart/internal/utils"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type OrderStore struct {
	db  *pgxpool.Pool
	log *zap.Logger
}

func NewOrderStore(db *pgxpool.Pool, log *zap.Logger) OrderStore {
	return OrderStore{db: db, log: log}
}

func (o *OrderStore) CreateOrder(orderID string, userID int) (bool, error) {
	var statusID int
	sql := `SELECT id FROM statuses WHERE status = 'new'`
	err := o.db.QueryRow(context.Background(), sql).Scan(&statusID)
	if err != nil {
		return false, fmt.Errorf("failed to get status ID for 'new': %w", err)
	}

	sql = `INSERT INTO orders (order_ID, user_ID, status_ID) VALUES ($1, $2, $3)`
	_, err = o.db.Exec(context.Background(), sql, orderID, userID, statusID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (o *OrderStore) GetOrderByID(orderID string) (*model.Order, error) {
	sql := `SELECT order_ID, user_ID, accrual, created_at, updated_at FROM orders WHERE order_ID = $1`
	row := o.db.QueryRow(context.Background(), sql, orderID)
	var order model.Order
	err := row.Scan(&order.OrderID, &order.UserID, &order.Accrual, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	order.Accrual = utils.Round(order.Accrual, 5)

	return &order, nil
}

func (o *OrderStore) UpdateSumAndStatusOrder(orderID string, status string, sum float64) (bool, error) {
	tx, err := o.db.Begin(context.Background())
	if err != nil {
		return false, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	sql := `SELECT id FROM statuses WHERE status = $1`
	row := tx.QueryRow(context.Background(), sql, status)
	var statusID int
	err = row.Scan(&statusID)
	if err != nil {
		return false, fmt.Errorf("error fetching status ID: %w", err)
	}

	sql = `SELECT user_id FROM orders WHERE order_id = $1`
	row = tx.QueryRow(context.Background(), sql, orderID)
	var userID int
	err = row.Scan(&userID)
	if err != nil {
		return false, fmt.Errorf("error fetching user ID: %w", err)
	}

	roundedSum := utils.Round(sum, 5)

	sql = `UPDATE orders SET status_ID = $1, accrual = $2, updated_at = NOW() WHERE order_id = $3`
	cmdTag, err := tx.Exec(context.Background(), sql, statusID, roundedSum, orderID)
	if err != nil {
		return false, fmt.Errorf("error updating order: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return false, fmt.Errorf("no order found with id: %s", orderID)
	}

	sql = `UPDATE users SET account = account + $1 WHERE id = $2`
	cmdTag, err = tx.Exec(context.Background(), sql, roundedSum, userID)
	if err != nil {
		return false, fmt.Errorf("error updating user account: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return false, fmt.Errorf("no user found with id: %d", userID)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return false, fmt.Errorf("error committing transaction: %w", err)
	}

	return true, nil
}

func (o *OrderStore) GetAllOrdersByUser(userID int64) ([]model.Order, error) {
	sql := `
		SELECT o.order_id, o.user_id, s.status, o.accrual, o.created_at, o.updated_at
		FROM orders o 
		LEFT JOIN statuses s ON o.status_id = s.id
		WHERE o.user_ID = $1
	`

	rows, err := o.db.Query(context.Background(), sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.OrderID, &order.UserID, &order.Status, &order.Accrual, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			return nil, fmt.Errorf("ошибка загрузки временной зоны: %w", err)
		}

		order.CreatedAt = order.CreatedAt.In(loc)

		order.UpdatedAt = order.UpdatedAt.In(loc)

		order.Status = strings.ToUpper(order.Status)

		order.Accrual = utils.Round(order.Accrual, 5)

		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return nil, nil
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OrderStore) GetAllWithdrawnByUser(userID int64) (float64, error) {
	var totalChange float64

	sql := `
		SELECT SUM(change) as total_change
		FROM balance_history
		WHERE user_id = $1
	`

	err := o.db.QueryRow(context.Background(), sql, userID).Scan(&totalChange)
	if err != nil {
		return 0, nil
	}

	totalChange = utils.Round(totalChange, 5)

	return totalChange, nil
}

func (o *OrderStore) GetAllUnhandlerOrders(unhandledStatus1, unhandledStatus2 int) ([]string, error) {
	sql := `SELECT order_id FROM orders WHERE status_id = $1 OR status_id = $2`

	rows, err := o.db.Query(context.Background(), sql, unhandledStatus1, unhandledStatus2)
	if err != nil {
		o.log.Error("error GetAllUnhandlerOrders", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var orders []string
	for rows.Next() {
		var order string
		err := rows.Scan(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
