package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/space-wanderer/microservices/order/internal/app"
	"github.com/space-wanderer/microservices/order/internal/config"
	"github.com/space-wanderer/microservices/platform/pkg/closer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

const configPath = "deploy/compose/order/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Инициализируем OpenTelemetry для отправки метрик в OTEL Collector
	ctx := context.Background()

	// Создаем ресурс с информацией о сервисе
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("order-service"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create resource: %v", err))
	}

	// Создаем OTLP exporter для отправки метрик в OTEL Collector
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create OTLP exporter: %v", err))
	}

	// Настраиваем OpenTelemetry meter provider
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter)),
	)
	otel.SetMeterProvider(meterProvider)

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, fmt.Sprintf("failed to create app order: %v", err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, fmt.Sprintf("failed to run app order: %v", err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, fmt.Sprintf("failed to close all: %v", err))
	}
}
