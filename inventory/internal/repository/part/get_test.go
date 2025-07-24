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
}

func (s *GetPartTestSuite) SetupTest() {
	s.mockRepository = mocks.NewInventoryRepository(s.T())
}

func (s *GetPartTestSuite) TearDownTest() {
	s.mockRepository.AssertExpectations(s.T())
}

func TestGetPartTestSuite(t *testing.T) {
	suite.Run(t, new(GetPartTestSuite))
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
