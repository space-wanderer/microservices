package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/space-wanderer/microservices/inventory/internal/model"
	"github.com/space-wanderer/microservices/inventory/internal/repository/mocks"
	repoModel "github.com/space-wanderer/microservices/inventory/internal/repository/model"
)

func TestService_ListParts(t *testing.T) {
	tests := []struct {
		name           string
		filter         *model.PartsFilter
		setupMock      func(*mocks.InventoryRepository)
		expectedResult []*model.Part
		expectedError  bool
	}{
		{
			name: "Успешное получение списка деталей",
			filter: &model.PartsFilter{
				Categories: []model.Category{model.CategoryEngine},
				Tags:       []string{"test"},
			},
			setupMock: func(mockRepo *mocks.InventoryRepository) {
				expectedFilter := &repoModel.PartsFilter{
					Categories: []repoModel.Category{repoModel.CategoryEngine},
					Tags:       []string{"test"},
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

				mockRepo.EXPECT().
					ListParts(mock.Anything, expectedFilter).
					Return(expectedParts, nil).
					Once()
			},
			expectedResult: []*model.Part{
				{
					UUID:          "test-uuid-1",
					Name:          "Test Engine 1",
					Description:   "Test engine description 1",
					Price:         1000.0,
					StockQuantity: 10,
					Category:      model.CategoryEngine,
				},
				{
					UUID:          "test-uuid-2",
					Name:          "Test Engine 2",
					Description:   "Test engine description 2",
					Price:         2000.0,
					StockQuantity: 20,
					Category:      model.CategoryEngine,
				},
			},
			expectedError: false,
		},
		{
			name: "Пустой список деталей",
			filter: &model.PartsFilter{
				Categories: []model.Category{model.CategoryFuel},
			},
			setupMock: func(mockRepo *mocks.InventoryRepository) {
				expectedFilter := &repoModel.PartsFilter{
					Categories: []repoModel.Category{repoModel.CategoryFuel},
				}

				mockRepo.EXPECT().
					ListParts(mock.Anything, expectedFilter).
					Return([]*repoModel.Part{}, nil).
					Once()
			},
			expectedResult: []*model.Part{},
			expectedError:  false,
		},
		{
			name: "Ошибка при получении списка",
			filter: &model.PartsFilter{
				Uuids: []string{"test-uuid-error"},
			},
			setupMock: func(mockRepo *mocks.InventoryRepository) {
				expectedFilter := &repoModel.PartsFilter{
					Uuids: []string{"test-uuid-error"},
				}

				mockRepo.EXPECT().
					ListParts(mock.Anything, expectedFilter).
					Return(nil, assert.AnError).
					Once()
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name:   "Фильтр без параметров",
			filter: &model.PartsFilter{},
			setupMock: func(mockRepo *mocks.InventoryRepository) {
				expectedFilter := &repoModel.PartsFilter{}

				mockRepo.EXPECT().
					ListParts(mock.Anything, expectedFilter).
					Return([]*repoModel.Part{}, nil).
					Once()
			},
			expectedResult: []*model.Part{},
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок репозитория
			mockRepo := mocks.NewInventoryRepository(t)

			// Настраиваем ожидания мока
			tt.setupMock(mockRepo)

			// Создаем сервис с моком
			service := NewService(mockRepo)

			// Выполняем тест
			result, err := service.ListParts(context.Background(), tt.filter)

			// Проверяем результаты
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tt.expectedResult))
			}
		})
	}
}

func TestService_ListParts_WithMock(t *testing.T) {
	// Создаем мок репозитория
	mockRepo := mocks.NewInventoryRepository(t)
	service := NewService(mockRepo)

	// Настраиваем ожидания мока
	expectedParts := []*repoModel.Part{
		{
			UUID:          "test-uuid-1",
			Name:          "Test Engine 1",
			Category:      repoModel.CategoryEngine,
			Price:         1000.0,
			StockQuantity: 10,
		},
		{
			UUID:          "test-uuid-2",
			Name:          "Test Engine 2",
			Category:      repoModel.CategoryEngine,
			Price:         2000.0,
			StockQuantity: 5,
		},
	}

	filter := &model.PartsFilter{
		Categories: []model.Category{model.CategoryEngine},
	}

	mockRepo.EXPECT().
		ListParts(mock.Anything, mock.Anything).
		Return(expectedParts, nil).
		Once()

	// Выполняем тест
	result, err := service.ListParts(context.Background(), filter)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "test-uuid-1", result[0].UUID)
	assert.Equal(t, "test-uuid-2", result[1].UUID)
}
