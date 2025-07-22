package part

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/space-wanderer/microservices/inventory/internal/repository/mocks"
	repoModel "github.com/space-wanderer/microservices/inventory/internal/repository/model"
)

type GetPartTestSuite struct {
	suite.Suite
	mockRepository *mocks.InventoryRepository
	repository     *repository
}

func (s *GetPartTestSuite) SetupTest() {
	s.mockRepository = mocks.NewInventoryRepository(s.T())
	s.repository = NewRepository()
}

func (s *GetPartTestSuite) TearDownTest() {
	s.mockRepository.AssertExpectations(s.T())
}

func TestGetPartTestSuite(t *testing.T) {
	suite.Run(t, new(GetPartTestSuite))
}

func (s *GetPartTestSuite) TestGetPart_Success() {
	// Arrange
	ctx := context.Background()
	uuid := "550e8400-e29b-41d4-a716-446655440001"

	// Act
	result, err := s.repository.GetPart(ctx, uuid)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), uuid, result.UUID)
	assert.Equal(s.T(), "Ионный двигатель X-2000", result.Name)
	assert.Equal(s.T(), repoModel.CategoryEngine, result.Category)
	assert.Equal(s.T(), 150000.0, result.Price)
	assert.Equal(s.T(), int64(5), result.StockQuantity)
}

func (s *GetPartTestSuite) TestGetPart_NotFound() {
	// Arrange
	ctx := context.Background()
	uuid := "non-existent-uuid"

	// Act
	result, err := s.repository.GetPart(ctx, uuid)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "part not found")
}

func (s *GetPartTestSuite) TestGetPart_EmptyUUID() {
	// Arrange
	ctx := context.Background()
	uuid := ""

	// Act
	result, err := s.repository.GetPart(ctx, uuid)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Contains(s.T(), err.Error(), "part not found")
}

func (s *GetPartTestSuite) TestGetPart_WithMockRepository() {
	// Arrange
	ctx := context.Background()
	uuid := "test-uuid-123"

	expectedPart := &repoModel.Part{
		UUID:          uuid,
		Name:          "Test Engine",
		Description:   "Test engine description",
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
			Name:    "Test Company",
			Country: "Test Country",
			Website: "https://test.com",
		},
		Tags:     []string{"test", "engine"},
		Metadata: map[string]*repoModel.Value{},
	}

	s.mockRepository.On("GetPart", ctx, uuid).Return(expectedPart, nil).Once()

	// Act - здесь мы тестируем мок, а не реальную логику
	result, err := s.mockRepository.GetPart(ctx, uuid)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), uuid, result.UUID)
	assert.Equal(s.T(), "Test Engine", result.Name)
}

func (s *GetPartTestSuite) TestGetPart_MockRepositoryError() {
	// Arrange
	ctx := context.Background()
	uuid := "error-uuid"
	expectedError := errors.New("database error")

	s.mockRepository.On("GetPart", ctx, uuid).Return(nil, expectedError).Once()

	// Act
	result, err := s.mockRepository.GetPart(ctx, uuid)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), expectedError, err)
}

func (s *GetPartTestSuite) TestGetPart_AllCategories() {
	// Arrange
	testCases := []struct {
		uuid     string
		name     string
		category repoModel.Category
		price    float64
	}{
		{
			uuid:     "550e8400-e29b-41d4-a716-446655440001",
			name:     "Ионный двигатель X-2000",
			category: repoModel.CategoryEngine,
			price:    150000.0,
		},
		{
			uuid:     "550e8400-e29b-41d4-a716-446655440003",
			name:     "Криогенное топливо H2-O2",
			category: repoModel.CategoryFuel,
			price:    50000.0,
		},
		{
			uuid:     "550e8400-e29b-41d4-a716-446655440005",
			name:     "Кварцевое окно QW-100",
			category: repoModel.CategoryPorthole,
			price:    25000.0,
		},
		{
			uuid:     "550e8400-e29b-41d4-a716-446655440007",
			name:     "Солнечная панель SP-500",
			category: repoModel.CategoryWing,
			price:    75000.0,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Act
			result, err := s.repository.GetPart(context.Background(), tc.uuid)

			// Assert
			assert.NoError(s.T(), err)
			assert.NotNil(s.T(), result)
			assert.Equal(s.T(), tc.uuid, result.UUID)
			assert.Equal(s.T(), tc.name, result.Name)
			assert.Equal(s.T(), tc.category, result.Category)
			assert.Equal(s.T(), tc.price, result.Price)
		})
	}
}

func (s *GetPartTestSuite) TestGetPart_ConcurrentAccess() {
	// Arrange
	ctx := context.Background()
	uuid := "550e8400-e29b-41d4-a716-446655440001"

	// Создаем каналы для синхронизации
	done := make(chan bool, 10)

	// Act - запускаем несколько горутин одновременно
	for i := 0; i < 10; i++ {
		go func() {
			result, err := s.repository.GetPart(ctx, uuid)
			assert.NoError(s.T(), err)
			assert.NotNil(s.T(), result)
			assert.Equal(s.T(), uuid, result.UUID)
			done <- true
		}()
	}

	// Ждем завершения всех горутин
	for i := 0; i < 10; i++ {
		<-done
	}
}

func (s *GetPartTestSuite) TestGetPart_WithManufacturer() {
	// Arrange
	ctx := context.Background()
	uuid := "550e8400-e29b-41d4-a716-446655440001"

	// Act
	result, err := s.repository.GetPart(ctx, uuid)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.NotNil(s.T(), result.Manufacturer)
	assert.Equal(s.T(), "КосмоТех", result.Manufacturer.Name)
	assert.Equal(s.T(), "Россия", result.Manufacturer.Country)
	assert.Equal(s.T(), "https://cosmotech.ru", result.Manufacturer.Website)
}

func (s *GetPartTestSuite) TestGetPart_WithDimensions() {
	// Arrange
	ctx := context.Background()
	uuid := "550e8400-e29b-41d4-a716-446655440001"

	// Act
	result, err := s.repository.GetPart(ctx, uuid)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.NotNil(s.T(), result.Dimensions)
	assert.Equal(s.T(), 120.0, result.Dimensions.Length)
	assert.Equal(s.T(), 80.0, result.Dimensions.Width)
	assert.Equal(s.T(), 60.0, result.Dimensions.Height)
	assert.Equal(s.T(), 250.0, result.Dimensions.Weight)
}

func (s *GetPartTestSuite) TestGetPart_WithTags() {
	// Arrange
	ctx := context.Background()
	uuid := "550e8400-e29b-41d4-a716-446655440001"

	// Act
	result, err := s.repository.GetPart(ctx, uuid)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.Tags, 4)
	assert.Contains(s.T(), result.Tags, "ионный")
	assert.Contains(s.T(), result.Tags, "двигатель")
	assert.Contains(s.T(), result.Tags, "межпланетный")
	assert.Contains(s.T(), result.Tags, "высокоэффективный")
}
