package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/space-wanderer/microservices/notification/internal/app"
	"github.com/space-wanderer/microservices/notification/internal/config"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

const configPath = "deploy/compose/notification/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	if err := config.Load(); err != nil {
		log.Fatalf("❌ Ошибка загрузки конфигурации: %v", err)
	}

	// Логгер инициализируется в app.go

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
