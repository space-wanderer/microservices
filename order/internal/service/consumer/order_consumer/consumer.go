package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/space-wanderer/microservices/order/internal/converter/kafka"
	orderService "github.com/space-wanderer/microservices/order/internal/service"
	"github.com/space-wanderer/microservices/platform/pkg/kafka"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

type Service struct {
	orderConsumer        kafka.Consumer
	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder
	orderService         orderService.OrderService
}

func NewService(orderConsumer kafka.Consumer, shipAssembledDecoder kafkaConverter.ShipAssembledDecoder, orderService orderService.OrderService) *Service {
	return &Service{
		orderConsumer:        orderConsumer,
		shipAssembledDecoder: shipAssembledDecoder,
		orderService:         orderService,
	}
}

func (s *Service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order consumer")

	err := s.orderConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "failed to consume order", zap.Error(err))
		return err
	}

	return nil
}
