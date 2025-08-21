package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/space-wanderer/microservices/notification/internal/converter/kafka"
	"github.com/space-wanderer/microservices/notification/internal/service/telegram"
	"github.com/space-wanderer/microservices/platform/pkg/kafka"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

type service struct {
	orderAssembledRecodeConsumer kafka.Consumer
	orderAssembledRecodeDecoder  kafkaConverter.ShipAssembledDecoder
	telegramService              *telegram.Service
}

func NewService(orderAssembledRecodeConsumer kafka.Consumer, orderAssembledRecodeDecoder kafkaConverter.ShipAssembledDecoder, telegramService *telegram.Service) *service {
	return &service{
		orderAssembledRecodeConsumer: orderAssembledRecodeConsumer,
		orderAssembledRecodeDecoder:  orderAssembledRecodeDecoder,
		telegramService:              telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting notificaiton consumer")

	err := s.orderAssembledRecodeConsumer.Consume(ctx, s.OrderAssembledHandler)
	if err != nil {
		logger.Error(ctx, "failed to consume order assembled", zap.Error(err))
		return err
	}

	return nil
}
