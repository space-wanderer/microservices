package order_consumer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/space-wanderer/microservices/assembly/internal/model"
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

	// Имитируем сборку корабля
	logger.Info(ctx, "🚀 Начинаем сборку корабля", zap.String("order_uuid", event.OrderUUID))

	// Используем таймер для сборки корабля
	timer := time.NewTimer(10 * time.Second)
	select {
	case <-timer.C:
		logger.Info(ctx, "✅ Корабль собран", zap.String("order_uuid", event.OrderUUID))
	case <-ctx.Done():
		timer.Stop()
		logger.Error(ctx, "❌ Сборка корабля прервана", zap.String("order_uuid", event.OrderUUID))
		return ctx.Err()
	}

	// Создаем событие ShipAssembledEvent
	shipAssembledEvent := model.ShipAssembledEvent{
		EventUUID:    uuid.New().String(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: 10,
	}

	// Отправляем событие через producer
	if err := s.producerService.ProduceShipAssembledEvent(ctx, shipAssembledEvent); err != nil {
		logger.Error(ctx, "❌ Ошибка отправки события ShipAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "📤 Событие ShipAssembled отправлено",
		zap.String("event_uuid", shipAssembledEvent.EventUUID),
		zap.String("order_uuid", shipAssembledEvent.OrderUUID),
	)

	return nil
}
