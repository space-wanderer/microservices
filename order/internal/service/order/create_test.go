package order

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/space-wanderer/microservices/order/internal/client/grpc/mocks"
	"github.com/space-wanderer/microservices/order/internal/model"
	repoMocks "github.com/space-wanderer/microservices/order/internal/repository/mocks"
	repoModel "github.com/space-wanderer/microservices/order/internal/repository/model"
)

type CreateOrderTestSuite struct {
	suite.Suite
	orderRepository *repoMocks.OrderRepository
	inventoryClient *mocks.InventoryClient
	paymentClient   *mocks.PaymentClient
	service         *service
}

func (s *CreateOrderTestSuite) SetupTest() {
	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = mocks.NewInventoryClient(s.T())
	s.paymentClient = mocks.NewPaymentClient(s.T())
	s.service = NewOrderService(s.orderRepository, s.inventoryClient, s.paymentClient, nil)
}

func (s *CreateOrderTestSuite) TearDownTest() {
	s.orderRepository.AssertExpectations(s.T())
	s.inventoryClient.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

func TestCreateOrderTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderTestSuite))
}

func (s *CreateOrderTestSuite) TestCreateOrder_Success() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"

	req := model.Order{
		OrderUUID:       "",
		UserUUID:        "550e8400-e29b-41d4-a716-446655440001",
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      0,
		TransactionUUID: nil,
		PaymentMethod:   model.PaymentMethodCard,
		Status:          model.StatusPendingPayment,
	}

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       "",
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      0,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPendingPayment,
	}

	expectedPart := &model.Part{
		UUID:  "550e8400-e29b-41d4-a716-446655440002",
		Name:  "Test Part",
		Price: 150.5,
	}

	expectedResult := model.Order{
		OrderUUID:  orderUUID,
		TotalPrice: float32(expectedPart.Price),
	}

	s.orderRepository.On("CreateOrder", ctx, expectedRepoOrder).Return(orderUUID, nil)
	s.inventoryClient.On("ListParts", mock.AnythingOfType("*context.timerCtx"), model.PartsFilter{
		Uuids: []string{"550e8400-e29b-41d4-a716-446655440002"},
	}).Return([]*model.Part{expectedPart}, nil)

	// Act
	result, err := s.service.CreateOrder(ctx, req)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResult, result)
}

func (s *CreateOrderTestSuite) TestCreateOrder_RepositoryError() {
	// Arrange
	ctx := context.Background()
	expectedError := errors.New("database error")

	req := model.Order{
		OrderUUID:       "",
		UserUUID:        "550e8400-e29b-41d4-a716-446655440001",
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      0,
		TransactionUUID: nil,
		PaymentMethod:   model.PaymentMethodCard,
		Status:          model.StatusPendingPayment,
	}

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       "",
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      0,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPendingPayment,
	}

	s.orderRepository.On("CreateOrder", ctx, expectedRepoOrder).Return("", expectedError)

	// Act
	result, err := s.service.CreateOrder(ctx, req)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *CreateOrderTestSuite) TestCreateOrder_InventoryClientError() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	expectedError := errors.New("inventory service error")

	req := model.Order{
		OrderUUID:       "",
		UserUUID:        "550e8400-e29b-41d4-a716-446655440001",
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      0,
		TransactionUUID: nil,
		PaymentMethod:   model.PaymentMethodCard,
		Status:          model.StatusPendingPayment,
	}

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       "",
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      0,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPendingPayment,
	}

	s.orderRepository.On("CreateOrder", ctx, expectedRepoOrder).Return(orderUUID, nil)
	s.inventoryClient.On("ListParts", mock.AnythingOfType("*context.timerCtx"), model.PartsFilter{
		Uuids: []string{"550e8400-e29b-41d4-a716-446655440002"},
	}).Return([]*model.Part(nil), expectedError)

	// Act
	result, err := s.service.CreateOrder(ctx, req)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "ошибка при получении информации о деталях")
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *CreateOrderTestSuite) TestCreateOrder_PartNotFound() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"

	req := model.Order{
		OrderUUID:       "",
		UserUUID:        "550e8400-e29b-41d4-a716-446655440001",
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      0,
		TransactionUUID: nil,
		PaymentMethod:   model.PaymentMethodCard,
		Status:          model.StatusPendingPayment,
	}

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       "",
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      0,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPendingPayment,
	}

	s.orderRepository.On("CreateOrder", ctx, expectedRepoOrder).Return(orderUUID, nil)
	s.inventoryClient.On("ListParts", mock.AnythingOfType("*context.timerCtx"), model.PartsFilter{
		Uuids: []string{"550e8400-e29b-41d4-a716-446655440002"},
	}).Return([]*model.Part{}, nil)

	// Act
	result, err := s.service.CreateOrder(ctx, req)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrOrderNotFound, err)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *CreateOrderTestSuite) TestCreateOrder_MultipleParts() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"

	req := model.Order{
		OrderUUID: "",
		UserUUID:  "550e8400-e29b-41d4-a716-446655440001",
		PartUuids: []string{
			"550e8400-e29b-41d4-a716-446655440002",
			"550e8400-e29b-41d4-a716-446655440003",
		},
		TotalPrice:      0,
		TransactionUUID: nil,
		PaymentMethod:   model.PaymentMethodCard,
		Status:          model.StatusPendingPayment,
	}

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       "",
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      0,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPendingPayment,
	}

	expectedParts := []*model.Part{
		{
			UUID:  "550e8400-e29b-41d4-a716-446655440002",
			Name:  "Test Part 1",
			Price: 150.5,
		},
		{
			UUID:  "550e8400-e29b-41d4-a716-446655440003",
			Name:  "Test Part 2",
			Price: 250.75,
		},
	}

	expectedTotalPrice := float32(expectedParts[0].Price + expectedParts[1].Price)

	expectedResult := model.Order{
		OrderUUID:  orderUUID,
		TotalPrice: expectedTotalPrice,
	}

	s.orderRepository.On("CreateOrder", ctx, expectedRepoOrder).Return(orderUUID, nil)

	// Первый вызов для первой детали
	s.inventoryClient.On("ListParts", mock.AnythingOfType("*context.timerCtx"), model.PartsFilter{
		Uuids: []string{"550e8400-e29b-41d4-a716-446655440002"},
	}).Return([]*model.Part{expectedParts[0]}, nil)

	// Второй вызов для второй детали
	s.inventoryClient.On("ListParts", mock.AnythingOfType("*context.timerCtx"), model.PartsFilter{
		Uuids: []string{"550e8400-e29b-41d4-a716-446655440003"},
	}).Return([]*model.Part{expectedParts[1]}, nil)

	// Act
	result, err := s.service.CreateOrder(ctx, req)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedResult, result)
	assert.Equal(s.T(), expectedTotalPrice, result.TotalPrice)
}
