package postgresql

import (
	"context"

	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderStore struct {
	db *pgxpool.Pool
}

func NewOrderStore(db *pgxpool.Pool) OrderStore {
	return OrderStore{db: db}
}

func (o *OrderStore) CreateOrder(order *model.Order) (bool, error) {
	sql := `INSERT INTO "order" (order_id, user_id, status_id, accrual, processed_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := o.db.Exec(context.Background(), sql, order.OrderID, order.UserID, order.StatusID, order.Accrual, order.ProcessedAt)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (o *OrderStore) UpdateStatusOrder(orderID int64, statusID int) (bool, error) {
	sql := `UPDATE "order" SET status_id = $1, updated_at = NOW() WHERE order_id = $2`
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
	sql := `SELECT order_id, user_id, status_id, accrual, processed_at, created_at, updated_at FROM "order" WHERE user_id = $1`
	rows, err := o.db.Query(context.Background(), sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.OrderID, &order.UserID, &order.StatusID, &order.Accrual, &order.ProcessedAt, &order.CreatedAt, &order.UpdatedAt)
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
