package order

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/space-wanderer/microservices/order/internal/repository/model"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

func (r *repository) UpdateOrder(ctx context.Context, order *model.Order) error {
	logger.Info(ctx, "🔄 Starting UpdateOrder",
		zap.String("order_uuid", order.OrderUUID),
		zap.String("status", string(order.Status)),
		zap.Any("transaction_uuid", order.TransactionUUID))

	conn, err := r.db.Acquire(ctx)
	if err != nil {
		logger.Error(ctx, "❌ Failed to acquire connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		logger.Error(ctx, "❌ Failed to begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			logger.Error(ctx, "❌ Failed to rollback transaction", zap.Error(rollbackErr))
			_ = rollbackErr
		}
	}()

	logger.Info(ctx, "💾 Executing UPDATE query",
		zap.String("order_uuid", order.OrderUUID),
		zap.String("status", string(order.Status)),
		zap.Any("transaction_uuid", order.TransactionUUID))

	result, err := tx.Exec(ctx, `
		UPDATE orders 
		SET user_uuid = $1, part_uuids = $2, total_price = $3, transaction_uuid = $4, payment_method = $5, status = $6, updated_at = NOW()
		WHERE order_uuid = $7
	`, order.UserUUID, order.PartUuids, order.TotalPrice, order.TransactionUUID, order.PaymentMethod, order.Status, order.OrderUUID)
	if err != nil {
		logger.Error(ctx, "❌ Failed to execute UPDATE query", zap.Error(err))
		return err
	}

	// Проверяем, что строка была обновлена
	rowsAffected := result.RowsAffected()
	logger.Info(ctx, "📊 UPDATE result",
		zap.String("order_uuid", order.OrderUUID),
		zap.Int64("rows_affected", rowsAffected))

	if rowsAffected == 0 {
		logger.Error(ctx, "❌ No rows affected", zap.String("order_uuid", order.OrderUUID))
		return pgx.ErrNoRows
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Error(ctx, "❌ Failed to commit transaction", zap.Error(err))
		return err
	}

	logger.Info(ctx, "✅ Successfully updated order",
		zap.String("order_uuid", order.OrderUUID),
		zap.String("status", string(order.Status)))

	return nil
}
