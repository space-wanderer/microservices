package app

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/space-wanderer/microservices/assembly/internal/config"
	"github.com/space-wanderer/microservices/platform/pkg/closer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	logger.Info(ctx, "🚀 Assembly Service started")

	// Запускаем consumer в горутине
	go func() {
		if err := a.diContainer.ConsumerService(ctx).RunConsumer(ctx); err != nil {
			logger.Error(ctx, "Consumer service error", zap.Error(err))
		}
	}()

	// Ждем завершения контекста
	<-ctx.Done()

	logger.Info(ctx, "🛑 Assembly Service shutting down...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return closer.CloseAll(shutdownCtx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initDi,
		a.initLogger,
		a.initCloser,
	}
	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(ctx context.Context) error {
	return config.Load()
}

func (a *App) initDi(ctx context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(ctx context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}
