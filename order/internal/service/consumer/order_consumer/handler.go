package order_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/space-wanderer/microservices/order/internal/model"
	"github.com/space-wanderer/microservices/platform/pkg/kafka/consumer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

func (s *Service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	event := s.shipAssembledDecoder.Decode(msg.Value)

	logger.Info(ctx, "Processing ShipAssembled message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int("build_time_sec", event.BuildTimeSec),
	)

	// Обновляем статус заказа на ASSEMBLED
	err := s.orderService.UpdateOrderStatus(ctx, event.OrderUUID, model.StatusAssembled)
	if err != nil {
		logger.Error(ctx, "Failed to update order status to ASSEMBLED",
			zap.String("order_uuid", event.OrderUUID),
			zap.Error(err))
		return err
	}

	logger.Info(ctx, "✅ Order status updated to ASSEMBLED-1",
		zap.String("order_uuid", event.OrderUUID))

	return nil
}
