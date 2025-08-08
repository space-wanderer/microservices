package integration

import (
	"context"

	"github.com/space-wanderer/microservices/platform/pkg/logger"
	"go.uber.org/zap"
)

func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	log := logger.Logger()
	log.Info(ctx, "Освобождение ресурсов тестового окружения...")

	cleanupTestEnvironment(ctx, env)

	log.Info(ctx, "Тестовое окружение освобождено")
}

// cleanupTestEnvironment — вспомогательная функция для освобождения ресурсов
func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер приложения", zap.Error(err))
		} else {
			logger.Info(ctx, "Контейнер приложения остановлен")
		}
	}

	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер MongoDB", zap.Error(err))
		} else {
			logger.Info(ctx, "Контейнер MongoDB остановлен")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "не удалось удалить сеть", zap.Error(err))
		} else {
			logger.Info(ctx, "Сеть удалена")
		}
	}
}
