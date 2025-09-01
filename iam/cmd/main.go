package main

import (
	"context"
	"fmt"
	"log"
	"syscall"

	"github.com/space-wanderer/microservices/platform/pkg/closer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"

	"github.com/space-wanderer/microservices/iam/internal/app"
	"github.com/space-wanderer/microservices/iam/internal/config"
)

const configPath = "deploy/compose/iam/.env"

func main() {
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("❌ Failed to load config: %v", err))
	}

	logger.Info(ctx, "🔧 Configuration loaded successfully")

	application, err := app.NewApp()
	if err != nil {
		panic(fmt.Sprintf("❌ Failed to create application: %v", err))
	}

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем приложение
	if err := application.Run(ctx); err != nil {
		log.Printf("❌ Application failed: %v", err)
	}
}
