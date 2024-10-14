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

func (o *OrderStore) CreateOrder(orderId int, userId int) (bool, error) {
	sql := `INSERT INTO orders (order_id, user_id, status_id) VALUES ($1, $2, $3)`
	_, err := o.db.Exec(context.Background(), sql, orderId, userId, 1)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (o *OrderStore) GetOrderById(orderId int) (*model.Order, error) {
	sql := `SELECT order_id, user_id, accrual, processed_at, created_at, updated_at FROM orders WHERE order_id = $1`
	row := o.db.QueryRow(context.Background(), sql, orderId)
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

func (o *OrderStore) UpdateStatusOrder(orderID int64, statusID int) (bool, error) {
	sql := `UPDATE orders SET status_id = $1, updated_at = NOW() WHERE order_id = $2`
	cmdTag, err := o.db.Exec(context.Background(), sql, statusID, orderID)
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
		SELECT o.order_id, o.user_id, s.status, o.accrual, o.processed_at, o.created_at, o.updated_at
		FROM orders o 
		LEFT JOIN statuses s ON o.status_id = s.id
		WHERE o.user_id = $1
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
		WHERE user_id = $1
	`

	err := o.db.QueryRow(context.Background(), sql, userID).Scan(&totalChange)
	if err != nil {
		return 0, nil
	}

	return 0, nil
}
