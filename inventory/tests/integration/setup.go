package integration

import (
	"context"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/space-wanderer/microservices/platform/pkg/logger"
	"github.com/space-wanderer/microservices/platform/pkg/testcontainers/app"
	"github.com/space-wanderer/microservices/platform/pkg/testcontainers/mongo"
	"github.com/space-wanderer/microservices/platform/pkg/testcontainers/network"
	testcontainers "github.com/space-wanderer/microservices/platform/pkg/testcontainers/path"
)

const (
	// Параметры для контейнеров
	inventoryAppName    = "inventory-app"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	// Переменные окружения приложения
	grpcPortKey = "GRPC_PORT"

	// Значения переменных окружения
	loggerLevelValue = "debug"
	startupTimeout   = 5 * time.Minute
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

// setupTestEnvironment — подготавливает тестовое окружение: сеть, контейнеры и возвращает структуру с ресурсами

func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "Подготавливаем тестовое окружение")

	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "Не удалось создать сеть", zap.Error(err))
	}
	logger.Info(ctx, "Сеть успешно создана")

	// Получаем переменные окружения для MongoDB с проверкой на наличие
	mongoUsername := getEnvWithLogging(ctx, "MONGO_INITDB_ROOT_USERNAME")
	mongoPassword := getEnvWithLogging(ctx, "MONGO_INITDB_ROOT_PASSWORD")
	mongoImageName := getEnvWithLogging(ctx, "MONGO_IMAGE_NAME")
	mongoDatabase := getEnvWithLogging(ctx, "MONGO_DATABASE")

	grpcPort := getEnvWithLogging(ctx, grpcPortKey)

	// Шаг 2: Запускаем контейнер с MongoDB
	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "Не удалось запустить контейнер с MongoDB", zap.Error(err))
	}

	logger.Info(ctx, "Контейнер с MongoDB запущен")

	// Шаг 3: Запускаем контейнер с приложением
	appEnv := map[string]string{
		// Переопределяем хост MongoDB для подключения к контейнеру из testcontainers
		"MONGO_HOST": generatedMongo.Config().Host,
		"MONGO_PORT": generatedMongo.Config().Port,
	}

	// Создаем настраиваемую стратегию ожидания с увеличенным таймаутом
	waitStrategy := wait.ForListeningPort(nat.Port(grpcPort + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile("deploy/docker/inventory", "Dockerfile"),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "Не удалось запустить контейнер с приложением", zap.Error(err))
	}

	logger.Info(ctx, "Контейнер с приложением запущен")

	logger.Info(ctx, "Тестовое окружение готово")
	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}
}

func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}
	return value
}
