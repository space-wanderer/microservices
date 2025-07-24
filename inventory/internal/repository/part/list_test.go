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

type ListPartsTestSuite struct {
	suite.Suite
	mockRepository *mocks.InventoryRepository
}

func (s *ListPartsTestSuite) SetupTest() {
	s.mockRepository = mocks.NewInventoryRepository(s.T())
}

func (s *ListPartsTestSuite) TearDownTest() {
	s.mockRepository.AssertExpectations(s.T())
}

func TestListPartsTestSuite(t *testing.T) {
	suite.Run(t, new(ListPartsTestSuite))
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
		Categories: []repoModel.Category{repoModel.CategoryEngine},
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
