package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/space-wanderer/microservices/platform/pkg/kafka/consumer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg consumer.Message) error {
	event := s.orderPaidRecodeDecoder.Decode(msg.Value)

	logger.Info(ctx, "Processing OrderPaid message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUUID),
	)

	// Отправляем уведомление в Telegram
	if err := s.telegramService.SendOrderPaidNotification(ctx, event.OrderUUID, event); err != nil {
		logger.Error(ctx, "Ошибка отправки уведомления OrderPaid", zap.Error(err))
		return err
	}

	return nil
}
