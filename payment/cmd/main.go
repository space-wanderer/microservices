package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/space-wanderer/microservices/payment/internal/app"
	"github.com/space-wanderer/microservices/payment/internal/config"
	"github.com/space-wanderer/microservices/platform/pkg/closer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

const configPath = "deploy/compose/payment/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, fmt.Sprintf("failed to create app: %v", err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, fmt.Sprintf("failed to run app: %v", err))
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
