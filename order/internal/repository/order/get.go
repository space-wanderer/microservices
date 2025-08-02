package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/space-wanderer/microservices/order/internal/repository/model"
)

func (r *repository) GetOrderByUuid(ctx context.Context, uuid string) (*model.Order, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var order model.Order
	err = conn.QueryRow(ctx, `
		SELECT order_uuid, user_uuid, part_uuids, total_price, payment_method, status
		FROM orders 
		WHERE order_uuid = $1
	`, uuid).Scan(&order.OrderUUID, &order.UserUUID, &order.PartUuids, &order.TotalPrice, &order.PaymentMethod, &order.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	return &order, nil
}
