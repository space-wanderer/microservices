package order

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/space-wanderer/microservices/order/internal/repository/model"
)

func (r *repository) UpdateOrder(ctx context.Context, order *model.Order) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			_ = rollbackErr
		}
	}()

	result, err := tx.Exec(ctx, `
		UPDATE orders 
		SET user_uuid = $1, part_uuids = $2, total_price = $3, payment_method = $4, status = $5
		WHERE order_uuid = $6
	`, order.UserUUID, order.PartUuids, order.TotalPrice, order.PaymentMethod, order.Status, order.OrderUUID)
	if err != nil {
		return err
	}

	// Проверяем, что строка была обновлена
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return tx.Commit(ctx)
}
