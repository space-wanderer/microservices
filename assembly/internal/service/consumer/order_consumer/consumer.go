package order_consumer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
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

	// МЕТРИКА: гистограмма длительности сборки
	assemblyDuration metric.Float64Histogram
}

func NewService(assemblyRecodeConsumer kafka.Consumer, assemblyRecodedDecoder kafkaConverter.AssemblyRecodedDecoder, producerService assemblyService.ProducerService) *service {
	// ИНИЦИАЛИЗАЦИЯ МЕТРИКИ: создаем гистограмму для измерения длительности сборки
	meter := otel.Meter("assembly-service")

	assemblyDuration, _ := meter.Float64Histogram(
		"assembly_duration_seconds",
		metric.WithDescription("Duration of assembly operations"),
		metric.WithUnit("s"),
	)

	return &service{
		assemblyRecodeConsumer: assemblyRecodeConsumer,
		assemblyRecodedDecoder: assemblyRecodedDecoder,
		producerService:        producerService,
		assemblyDuration:       assemblyDuration,
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
