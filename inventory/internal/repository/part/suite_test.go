package part

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/space-wanderer/microservices/inventory/internal/repository/mocks"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  Не удалось загрузить .env файл: %v", err)
	}
}

type RepositorySuite struct {
	suite.Suite
	mockRepository *mocks.InventoryRepository
	repository     *repository
	mongoClient    *mongo.Client
	testDB         *mongo.Database
}

func (s *RepositorySuite) SetupTest() {
	s.mockRepository = mocks.NewInventoryRepository(s.T())

	// Получаем конфигурацию для тестов
	mongoURI := os.Getenv("MONGODB_URI")
	// Подключаемся к тестовой MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		s.T().Logf("failed to connect to MongoDB: %v", err)
		s.T().FailNow()
	}

	// Проверяем подключение
	err = client.Ping(ctx, nil)
	if err != nil {
		s.T().Logf("failed to ping MongoDB: %v", err)
		s.T().FailNow()
	}

	s.mongoClient = client
	s.testDB = client.Database("inventory-test")

	// Создаем репозиторий с тестовой базой
	s.repository = NewRepository(s.testDB)
}

func (s *RepositorySuite) TearDownTest() {
	s.mockRepository.AssertExpectations(s.T())

	// Очищаем тестовую базу данных
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Удаляем коллекцию parts
	err := s.testDB.Collection("parts").Drop(ctx)
	if err != nil {
		s.T().Logf("failed to drop test collection: %v", err)
	}
}

func (s *RepositorySuite) TearDownSuite() {
	// Закрываем соединение с MongoDB
	if s.mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.mongoClient.Disconnect(ctx)
		if err != nil {
			s.T().Logf("failed to disconnect from MongoDB: %v", err)
		}
	}
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
