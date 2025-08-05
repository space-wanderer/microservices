package integration

import (
	"context"

	"github.com/space-wanderer/microservices/platform/pkg/logger"
	"go.uber.org/zap"
)

// teardownTestEnvironment — освобождает все ресурсы тестового окружения
func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	log := logger.Logger()
	log.Info(ctx, "Очистка тестового окружения...")

	cleanupTestEnvironment(ctx, env)

	log.Info(ctx, "✅ Тестовое окружение успешно очищено")
}

func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		if err := env.App.Terminate(ctx); err != nil {
			logger.Error(ctx, "Не удалось завершить контейнер с приложением", zap.Error(err))
		} else {
			logger.Info(ctx, "Контейнер с приложением остановлен")
		}
	}

	if env.Mongo != nil {
		if err := env.Mongo.Terminate(ctx); err != nil {
			logger.Error(ctx, "Не удалось завершить контейнер MongoDB", zap.Error(err))
		} else {
			logger.Info(ctx, "Контейнер MongoDB остановлен")
		}
	}

	if env.Network != nil {
		if err := env.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "Не удалось удалить сеть", zap.Error(err))
		} else {
			logger.Info(ctx, "Сеть удалена")
		}
	}
}
