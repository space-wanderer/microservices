package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/space-wanderer/microservices/platform/pkg/kafka/consumer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

func (s *service) OrderAssembledHandler(ctx context.Context, msg consumer.Message) error {
	event := s.orderAssembledRecodeDecoder.Decode(msg.Value)

	logger.Info(ctx, "Processing ShipAssembled message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Any("build_time_sec", event.BuildTimeSec),
	)

	// Отправляем уведомление в Telegram
	if err := s.telegramService.SendShipAssembledNotification(ctx, event.OrderUUID, event); err != nil {
		logger.Error(ctx, "Ошибка отправки уведомления ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
