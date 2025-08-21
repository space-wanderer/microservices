package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/space-wanderer/microservices/notification/internal/converter/kafka"
	"github.com/space-wanderer/microservices/notification/internal/service/telegram"
	"github.com/space-wanderer/microservices/platform/pkg/kafka"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

type service struct {
	orderPaidRecodeConsumer kafka.Consumer
	orderPaidRecodeDecoder  kafkaConverter.OrderPaidDecoder
	telegramService         *telegram.Service
}

func NewService(orderPaidRecodeConsumer kafka.Consumer, orderPaidRecodeDecoder kafkaConverter.OrderPaidDecoder, telegramService *telegram.Service) *service {
	return &service{
		orderPaidRecodeConsumer: orderPaidRecodeConsumer,
		orderPaidRecodeDecoder:  orderPaidRecodeDecoder,
		telegramService:         telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting notificaiton consumer")

	err := s.orderPaidRecodeConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "failed to consume order paid", zap.Error(err))
		return err
	}

	return nil
}
