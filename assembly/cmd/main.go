package main

import (
	"context"
	"log"
	"syscall"

	"github.com/space-wanderer/microservices/assembly/internal/app"
	"github.com/space-wanderer/microservices/platform/pkg/closer"
)

func main() {
	ctx := context.Background()

	// Создаем приложение
	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to create application: %v", err)
	}

	// Настраиваем обработку сигналов для graceful shutdown
	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	// Создаем контекст с отменой для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем приложение
	if err := application.Run(ctx); err != nil {
		log.Printf("❌ Application error: %v", err)
		// Не используем defer cancel() с os.Exit
		return
	}
}
