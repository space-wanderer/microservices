package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/space-wanderer/microservices/assembly/internal/converter/kafka"
	assemblyService "github.com/space-wanderer/microservices/assembly/internal/service"
	"github.com/space-wanderer/microservices/platform/pkg/kafka"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

type service struct {
	assemblyRecodeConsumer kafka.Consumer
	assemblyRecodedDecoder kafkaConverter.AssemblyRecodedDecoder
	producerService        assemblyService.ProducerService
}

func NewService(assemblyRecodeConsumer kafka.Consumer, assemblyRecodedDecoder kafkaConverter.AssemblyRecodedDecoder, producerService assemblyService.ProducerService) *service {
	return &service{
		assemblyRecodeConsumer: assemblyRecodeConsumer,
		assemblyRecodedDecoder: assemblyRecodedDecoder,
		producerService:        producerService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order consumer")

	err := s.assemblyRecodeConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "failed to consume order", zap.Error(err))
		return err
	}

	return nil
}
