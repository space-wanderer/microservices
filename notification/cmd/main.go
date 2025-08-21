package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/space-wanderer/microservices/notification/internal/app"
	"github.com/space-wanderer/microservices/notification/internal/config"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalf("❌ Ошибка загрузки конфигурации: %v", err)
	}

	if err := logger.Init(config.AppConfig().Logger.Level(), config.AppConfig().Logger.AsJson()); err != nil {
		log.Fatalf("❌ Ошибка инициализации логгера: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := app.NewApp()

	// Обрабатываем сигналы завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info(ctx, "🛑 Получен сигнал завершения")
		cancel()
	}()

	// Запускаем приложение
	if err := s.Run(ctx); err != nil {
		log.Printf("❌ Ошибка запуска приложения: %v", err)
		return
	}
}
