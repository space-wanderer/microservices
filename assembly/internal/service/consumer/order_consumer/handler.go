package order_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/space-wanderer/microservices/platform/pkg/kafka/consumer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	event := s.assemblyRecodedDecoder.Decode(msg.Value)

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUUID),
	)

	return nil
}
