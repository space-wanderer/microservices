package part

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/space-wanderer/microservices/inventory/internal/repository/mocks"
	repoModel "github.com/space-wanderer/microservices/inventory/internal/repository/model"
)

type ListPartsTestSuite struct {
	suite.Suite
	mockRepository *mocks.InventoryRepository
	repository     *repository
	mongoClient    *mongo.Client
	testDB         *mongo.Database
}

func (s *ListPartsTestSuite) SetupTest() {
	s.mockRepository = mocks.NewInventoryRepository(s.T())

	// Подключаемся к тестовой MongoDB с аутентификацией
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Используем те же параметры подключения, что и в main.go
	mongoURI := "mongodb://inventory-service-user:inventory-service-password@localhost:27017/inventory-test?authSource=admin"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		s.T().Fatalf("failed to connect to MongoDB: %v", err)
	}

	// Проверяем подключение
	err = client.Ping(ctx, nil)
	if err != nil {
		s.T().Fatalf("failed to ping MongoDB: %v", err)
	}

	s.mongoClient = client
	s.testDB = client.Database("inventory-test")

	// Создаем репозиторий с тестовой базой
	s.repository = NewRepository(s.testDB)
}

func (s *ListPartsTestSuite) TearDownTest() {
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

func (s *ListPartsTestSuite) TearDownSuite() {
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

func TestListPartsTestSuite(t *testing.T) {
	suite.Run(t, new(ListPartsTestSuite))
}

func (s *ListPartsTestSuite) TestListParts_NoFilter() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 8) // Всего 8 деталей в репозитории
}

func (s *ListPartsTestSuite) TestListParts_NilFilter() {
	// Arrange
	ctx := context.Background()

	// Act
	result, err := s.repository.ListParts(ctx, nil)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 8) // Всего 8 деталей в репозитории
}

func (s *ListPartsTestSuite) TestListParts_FilterByUUID() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Uuids: []string{"550e8400-e29b-41d4-a716-446655440001"},
	}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 1)
	assert.Equal(s.T(), "550e8400-e29b-41d4-a716-446655440001", result[0].UUID)
	assert.Equal(s.T(), "Ионный двигатель X-2000", result[0].Name)
}

func (s *ListPartsTestSuite) TestListParts_FilterByMultipleUUIDs() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Uuids: []string{
			"550e8400-e29b-41d4-a716-446655440001",
			"550e8400-e29b-41d4-a716-446655440002",
		},
	}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 2)

	// Проверяем, что получили нужные детали
	uuids := make([]string, len(result))
	for i, part := range result {
		uuids[i] = part.UUID
	}
	assert.Contains(s.T(), uuids, "550e8400-e29b-41d4-a716-446655440001")
	assert.Contains(s.T(), uuids, "550e8400-e29b-41d4-a716-446655440002")
}

func (s *ListPartsTestSuite) TestListParts_FilterByCategory() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Categories: []repoModel.Category{repoModel.CategoryEngine},
	}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 2) // 2 двигателя

	// Проверяем, что все детали имеют категорию ENGINE
	for _, part := range result {
		assert.Equal(s.T(), repoModel.CategoryEngine, part.Category)
	}
}

func (s *ListPartsTestSuite) TestListParts_FilterByMultipleCategories() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Categories: []repoModel.Category{
			repoModel.CategoryEngine,
			repoModel.CategoryFuel,
		},
	}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 4) // 2 двигателя + 2 топлива

	// Проверяем, что все детали имеют нужные категории
	for _, part := range result {
		assert.True(s.T(),
			part.Category == repoModel.CategoryEngine ||
				part.Category == repoModel.CategoryFuel)
	}
}

func (s *ListPartsTestSuite) TestListParts_FilterByManufacturerCountry() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		ManufacturerCountries: []string{"Россия"},
	}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)

	// Проверяем, что все детали произведены в России
	for _, part := range result {
		assert.NotNil(s.T(), part.Manufacturer)
		assert.Equal(s.T(), "Россия", part.Manufacturer.Country)
	}
}

func (s *ListPartsTestSuite) TestListParts_FilterByTags() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Tags: []string{"двигатель"},
	}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 2) // 2 двигателя

	// Проверяем, что все детали имеют тег "двигатель"
	for _, part := range result {
		assert.Contains(s.T(), part.Tags, "двигатель")
	}
}

func (s *ListPartsTestSuite) TestListParts_FilterByMultipleTags() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Tags: []string{"двигатель", "топливо"},
	}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)

	// Проверяем, что все детали имеют хотя бы один из тегов
	for _, part := range result {
		hasEngineTag := false
		hasFuelTag := false
		for _, tag := range part.Tags {
			if tag == "двигатель" {
				hasEngineTag = true
			}
			if tag == "топливо" {
				hasFuelTag = true
			}
		}
		assert.True(s.T(), hasEngineTag || hasFuelTag)
	}
}

func (s *ListPartsTestSuite) TestListParts_FilterByNames() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Names: []string{"Ионный двигатель X-2000"},
	}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 1)
	assert.Equal(s.T(), "Ионный двигатель X-2000", result[0].Name)
}

func (s *ListPartsTestSuite) TestListParts_ComplexFilter() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Categories:            []repoModel.Category{repoModel.CategoryEngine},
		ManufacturerCountries: []string{"Россия"},
		Tags:                  []string{"ионный"},
	}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 1) // Только ионный двигатель из России

	part := result[0]
	assert.Equal(s.T(), repoModel.CategoryEngine, part.Category)
	assert.Equal(s.T(), "Россия", part.Manufacturer.Country)
	assert.Contains(s.T(), part.Tags, "ионный")
}

func (s *ListPartsTestSuite) TestListParts_NoMatches() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Uuids: []string{"non-existent-uuid"},
	}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 0)
}

func (s *ListPartsTestSuite) TestListParts_EmptyFilter() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []repoModel.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 8) // Все детали
}

func (s *ListPartsTestSuite) TestListParts_FilterByNonExistentCategory() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Categories: []repoModel.Category{repoModel.CategoryUnknown},
	}

	// Act
	result, err := s.repository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 0) // Нет деталей с неизвестной категорией
}

func (s *ListPartsTestSuite) TestListParts_WithMockRepository() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Categories: []repoModel.Category{repoModel.CategoryEngine},
	}

	expectedParts := []*repoModel.Part{
		{
			UUID:          "test-uuid-1",
			Name:          "Test Engine 1",
			Description:   "Test engine description 1",
			Price:         1000.0,
			StockQuantity: 10,
			Category:      repoModel.CategoryEngine,
			Dimensions: &repoModel.Dimensions{
				Length: 10.0,
				Width:  5.0,
				Height: 2.0,
				Weight: 50.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Test Company 1",
				Country: "Test Country 1",
				Website: "https://test1.com",
			},
			Tags:     []string{"test", "engine"},
			Metadata: map[string]*repoModel.Value{},
		},
		{
			UUID:          "test-uuid-2",
			Name:          "Test Engine 2",
			Description:   "Test engine description 2",
			Price:         2000.0,
			StockQuantity: 20,
			Category:      repoModel.CategoryEngine,
			Dimensions: &repoModel.Dimensions{
				Length: 20.0,
				Width:  10.0,
				Height: 4.0,
				Weight: 100.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Test Company 2",
				Country: "Test Country 2",
				Website: "https://test2.com",
			},
			Tags:     []string{"test", "engine"},
			Metadata: map[string]*repoModel.Value{},
		},
	}

	s.mockRepository.On("ListParts", ctx, filter).Return(expectedParts, nil).Once()

	// Act
	result, err := s.mockRepository.ListParts(ctx, filter)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), "test-uuid-1", result[0].UUID)
	assert.Equal(s.T(), "test-uuid-2", result[1].UUID)
}

func (s *ListPartsTestSuite) TestListParts_MockRepositoryError() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Uuids: []string{"test-uuid-error"},
	}
	expectedError := errors.New("database error")

	s.mockRepository.On("ListParts", ctx, filter).Return(nil, expectedError).Once()

	// Act
	result, err := s.mockRepository.ListParts(ctx, filter)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedError, err)
}

func (s *ListPartsTestSuite) TestListParts_ConcurrentAccess() {
	// Arrange
	ctx := context.Background()
	filter := &repoModel.PartsFilter{
		Categories: []repoModel.Category{repoModel.CategoryEngine},
	}

	// Создаем каналы для синхронизации
	done := make(chan bool, 5)

	// Act - запускаем несколько горутин одновременно
	for i := 0; i < 5; i++ {
		go func() {
			result, err := s.repository.ListParts(ctx, filter)
			assert.NoError(s.T(), err)
			assert.NotNil(s.T(), result)
			assert.Len(s.T(), result, 2) // 2 двигателя
			done <- true
		}()
	}

	// Ждем завершения всех горутин
	for i := 0; i < 5; i++ {
		<-done
	}
}
