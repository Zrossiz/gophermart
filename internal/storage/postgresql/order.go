package postgresql

import (
	"context"

	"github.com/Zrossiz/gophermart/internal/model"
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

func (o *OrderStore) CreateOrder(orderID int, userID int) (bool, error) {
	sql := `INSERT INTO orders (order_ID, user_ID, status_ID) VALUES ($1, $2, $3)`
	_, err := o.db.Exec(context.Background(), sql, orderID, userID, 1)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (o *OrderStore) GetOrderByID(orderID int) (*model.Order, error) {
	sql := `SELECT order_ID, user_ID, accrual, processed_at, created_at, updated_at FROM orders WHERE order_ID = $1`
	row := o.db.QueryRow(context.Background(), sql, orderID)
	var order model.Order
	err := row.Scan(&order.OrderID, &order.UserID, &order.Accrual, &order.ProcessedAt, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}

func (o *OrderStore) UpdateSumAndStatusOrder(orderID int64, status string, sum float64) (bool, error) {
	sql := `SELECT ID FROM statuses WHERE status = $1`
	row := o.db.QueryRow(context.Background(), sql, status)
	var statusID int
	err := row.Scan(&statusID)
	if err != nil {
		return false, nil
	}

	sql = `UPDATE orders SET status_ID = $1, accrual = $2 updated_at = NOW() WHERE order_ID = $3`
	cmdTag, err := o.db.Exec(context.Background(), sql, statusID, sum, orderID)
	if err != nil {
		return false, err
	}
	if cmdTag.RowsAffected() == 0 {
		return false, nil
	}
	return true, nil
}

func (o *OrderStore) GetAllOrdersByUser(userID int64) ([]model.Order, error) {
	sql := `
		SELECT o.order_ID, o.user_ID, s.status, o.accrual, o.processed_at, o.created_at, o.updated_at
		FROM orders o 
		LEFT JOIN statuses s ON o.status_ID = s.ID
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
		err := rows.Scan(&order.OrderID, &order.UserID, &order.Status, &order.Accrual, &order.ProcessedAt, &order.CreatedAt, &order.UpdatedAt)
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

func (o *OrderStore) GetAllWithdrawnByUser(userID int64) (float64, error) {
	var totalChange float64

	sql := `
		SELECT SUM(change) as total_change
		FROM balance_history
		WHERE user_ID = $1
	`

	err := o.db.QueryRow(context.Background(), sql, userID).Scan(&totalChange)
	if err != nil {
		return 0, nil
	}

	return 0, nil
}

func (o *OrderStore) GetAllUnhandlerOrders(unhandledStatus1, unhandledStatus2 int) ([]model.Order, error) {
	sql := `SELECT order_ID, user_ID, status_ID FROM orders WHERE status_ID = $1 OR status_ID = $2`

	rows, err := o.db.Query(context.Background(), sql, unhandledStatus1, unhandledStatus2)
	if err != nil {
		o.log.Error("error GetAllWithdrawnByUser", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		for rows.Next() {
			var order model.Order
			err := rows.Scan(&order.OrderID, &order.UserID, &order.Status)
			if err != nil {
				return nil, err
			}
			orders = append(orders, order)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
