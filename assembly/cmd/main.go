package main

import (
	"context"
	"fmt"
	"log"
	"syscall"

	"github.com/space-wanderer/microservices/assembly/internal/app"
	"github.com/space-wanderer/microservices/assembly/internal/config"
	"github.com/space-wanderer/microservices/platform/pkg/closer"
)

const configPath = "deploy/compose/assembly/.env"

func main() {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Создаем приложение
	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to create application: %v", err)
	}

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем приложение
	if err := application.Run(ctx); err != nil {
		log.Printf("❌ Application error: %v", err)
		return
	}
}
