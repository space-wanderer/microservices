package app

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/space-wanderer/microservices/notification/internal/config"
	"github.com/space-wanderer/microservices/platform/pkg/closer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
	closer      *closer.Closer
}

func NewApp() *App {
	return &App{
		diContainer: NewDiContainer(),
		closer:      closer.New(),
	}
}

func (app *App) Run(ctx context.Context) error {
	// Инициализируем логгер с новой конфигурацией
	if err := app.initLogger(ctx); err != nil {
		return err
	}

	logger.Info(ctx, "🚀 Запуск Notification Service")

	// Запускаем OrderPaid consumer
	orderPaidConsumerService := app.diContainer.OrderPaidConsumerService(ctx)
	go func() {
		if err := orderPaidConsumerService.RunConsumer(ctx); err != nil {
			logger.Error(ctx, "OrderPaid consumer service error", zap.Error(err))
		}
	}()

	// Запускаем OrderAssembled consumer
	orderAssembledConsumerService := app.diContainer.OrderAssembledConsumerService(ctx)
	go func() {
		if err := orderAssembledConsumerService.RunConsumer(ctx); err != nil {
			logger.Error(ctx, "OrderAssembled consumer service error", zap.Error(err))
		}
	}()

	<-ctx.Done()

	logger.Info(ctx, "🛑 Получен сигнал завершения, начинаем graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := app.closer.CloseAll(shutdownCtx); err != nil {
		logger.Error(ctx, "Ошибка при закрытии ресурсов", zap.Error(err))
	}

	logger.Info(ctx, "✅ Notification Service завершен")
	return nil
}

func (app *App) initLogger(ctx context.Context) error {
	// Инициализируем логгер с новой конфигурацией
	return logger.InitWithConfig(config.AppConfig().Logger)
}
