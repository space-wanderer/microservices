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

func TestService_GetPart(t *testing.T) {
	tests := []struct {
		name           string
		uuid           string
		setupMock      func(*mocks.InventoryRepository)
		expectedResult *model.Part
		expectedError  bool
	}{
		{
			name: "Успешное получение детали",
			uuid: "test-uuid-123",
			setupMock: func(mockRepo *mocks.InventoryRepository) {
				expectedPart := &repoModel.Part{
					UUID:          "test-uuid-123",
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

				mockRepo.EXPECT().
					GetPart(mock.Anything, "test-uuid-123").
					Return(expectedPart, nil).
					Once()
			},
			expectedResult: &model.Part{
				UUID:          "test-uuid-123",
				Name:          "Test Engine",
				Description:   "Test engine description",
				Price:         1000.0,
				StockQuantity: 10,
				Category:      model.CategoryEngine,
			},
			expectedError: false,
		},
		{
			name: "Деталь не найдена",
			uuid: "not-found-uuid",
			setupMock: func(mockRepo *mocks.InventoryRepository) {
				mockRepo.EXPECT().
					GetPart(mock.Anything, "not-found-uuid").
					Return(nil, assert.AnError).
					Once()
			},
			expectedResult: nil,
			expectedError:  true,
		},
		{
			name: "Деталь с пустым UUID",
			uuid: "",
			setupMock: func(mockRepo *mocks.InventoryRepository) {
				mockRepo.EXPECT().
					GetPart(mock.Anything, "").
					Return(nil, assert.AnError).
					Once()
			},
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок репозитория
			mockRepo := newInventoryRepository(t)

			// Настраиваем ожидания мока
			tt.setupMock(mockRepo)

			// Создаем сервис с моком
			service := NewService(mockRepo)

			// Выполняем тест
			result, err := service.GetPart(context.Background(), tt.uuid)

			// Проверяем результаты
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.uuid, result.UUID)
			}
		})
	}
}
func newInventoryRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *mocks.InventoryRepository {
	mock := mocks.NewInventoryRepository(t)
	return mock
}
