package app

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/space-wanderer/microservices/iam/internal/config"
	"github.com/space-wanderer/microservices/platform/pkg/closer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
	authV1 "github.com/space-wanderer/microservices/shared/pkg/proto/auth/v1"
	userV1 "github.com/space-wanderer/microservices/shared/pkg/proto/user/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
}

func NewApp() (*App, error) {
	a := &App{
		diContainer: NewDiContainer(),
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	// Инициализируем логгер с новой конфигурацией
	if err := a.initLogger(ctx); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	logger.Info(ctx, "🚀 Starting IAM Service...")

	// Инициализируем gRPC сервер
	a.initGRPCServer(ctx)

	// Запускаем сервер
	grpcConfig := config.AppConfig().IamGrpc
	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", grpcConfig.Address(), err)
	}

	logger.Info(ctx, "🌐 gRPC server listening on", zap.String("address", grpcConfig.Address()))

	// Добавляем graceful shutdown
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		logger.Info(ctx, "🛑 Stopping gRPC server...")
		a.grpcServer.GracefulStop()
		return nil
	})

	// Запускаем сервер в горутине
	go func() {
		if err := a.grpcServer.Serve(lis); err != nil {
			logger.Error(ctx, "❌ gRPC server failed", zap.Error(err))
		}
	}()

	// Ждем сигнал завершения
	<-ctx.Done()
	logger.Info(ctx, "👋 IAM Service stopped")

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) {
	_ = ctx // Используем контекст для получения API
	// Создаем gRPC сервер
	a.grpcServer = grpc.NewServer()

	// Регистрируем сервисы
	userAPI := a.diContainer.UserV1API(ctx)
	authAPI := a.diContainer.AuthV1API(ctx)

	userV1.RegisterUserServiceServer(a.grpcServer, userAPI)
	authV1.RegisterAuthServiceServer(a.grpcServer, authAPI)

	// Добавляем health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(a.grpcServer, healthServer)

	// Устанавливаем статус сервисов
	healthServer.SetServingStatus("user.v1.UserService", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("auth.v1.AuthService", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(a.grpcServer)

	logger.Info(ctx, "✅ gRPC services registered successfully")
}

func (a *App) initLogger(ctx context.Context) error {
	// Инициализируем логгер с новой конфигурацией
	return logger.InitWithConfig(config.AppConfig().Logger)
}
