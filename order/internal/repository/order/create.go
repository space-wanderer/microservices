package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/space-wanderer/microservices/order/internal/repository/model"
)

func (r *repository) CreateOrder(ctx context.Context, req *model.Order) (string, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return "", err
	}

	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			_ = rollbackErr
		}
	}()

	orderUUID := uuid.New().String()

	_, err = tx.Exec(ctx, `
		INSERT INTO orders (order_uuid, user_uuid, part_uuids, total_price, payment_method, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, orderUUID, req.UserUUID, req.PartUuids, req.TotalPrice, req.PaymentMethod, req.Status)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return orderUUID, nil
}
