package order

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/space-wanderer/microservices/order/internal/client/grpc/mocks"
	"github.com/space-wanderer/microservices/order/internal/model"
	repoMocks "github.com/space-wanderer/microservices/order/internal/repository/mocks"
	repoModel "github.com/space-wanderer/microservices/order/internal/repository/model"
)

type CancelOrderTestSuite struct {
	suite.Suite
	orderRepository *repoMocks.OrderRepository
	inventoryClient *mocks.InventoryClient
	paymentClient   *mocks.PaymentClient
	service         *service
}

func (s *CancelOrderTestSuite) SetupTest() {
	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = mocks.NewInventoryClient(s.T())
	s.paymentClient = mocks.NewPaymentClient(s.T())
	s.service = NewOrderService(s.orderRepository, s.inventoryClient, s.paymentClient, nil)
}

func (s *CancelOrderTestSuite) TearDownTest() {
	s.orderRepository.AssertExpectations(s.T())
	s.inventoryClient.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

func TestCancelOrderTestSuite(t *testing.T) {
	suite.Run(t, new(CancelOrderTestSuite))
}

func (s *CancelOrderTestSuite) TestCancelOrderByUuid_Success() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "550e8400-e29b-41d4-a716-446655440001"

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPendingPayment,
	}

	expectedModelOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: nil,
		PaymentMethod:   model.PaymentMethodCard,
		Status:          model.StatusCanceled,
	}

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusCanceled,
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(repoOrder, nil)
	s.orderRepository.On("UpdateOrder", ctx, expectedRepoOrder).Return(nil)

	// Act
	result, err := s.service.CancelOrderByUuid(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedModelOrder, result)
}

func (s *CancelOrderTestSuite) TestCancelOrderByUuid_GetOrderError() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	expectedError := errors.New("database error")

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(nil, expectedError)

	// Act
	result, err := s.service.CancelOrderByUuid(ctx, orderUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrOrderNotFound, err)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *CancelOrderTestSuite) TestCancelOrderByUuid_OrderAlreadyPaid() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "550e8400-e29b-41d4-a716-446655440001"
	transactionUUID := "550e8400-e29b-41d4-a716-446655440003"

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPaid, // Уже оплачен
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(repoOrder, nil)

	// Act
	result, err := s.service.CancelOrderByUuid(ctx, orderUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrOrderCannotBeCancelled, err)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *CancelOrderTestSuite) TestCancelOrderByUuid_UpdateOrderError() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "550e8400-e29b-41d4-a716-446655440001"
	expectedError := errors.New("update error")

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPendingPayment,
	}

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusCanceled,
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(repoOrder, nil)
	s.orderRepository.On("UpdateOrder", ctx, expectedRepoOrder).Return(expectedError)

	// Act
	result, err := s.service.CancelOrderByUuid(ctx, orderUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to update order")
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *CancelOrderTestSuite) TestCancelOrderByUuid_AlreadyCanceled() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "550e8400-e29b-41d4-a716-446655440001"

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusCanceled, // Уже отменен
	}

	expectedModelOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: nil,
		PaymentMethod:   model.PaymentMethodCard,
		Status:          model.StatusCanceled,
	}

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusCanceled,
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(repoOrder, nil)
	s.orderRepository.On("UpdateOrder", ctx, expectedRepoOrder).Return(nil)

	// Act
	result, err := s.service.CancelOrderByUuid(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedModelOrder, result)
}

func (s *CancelOrderTestSuite) TestCancelOrderByUuid_EmptyUUID() {
	// Arrange
	ctx := context.Background()
	orderUUID := ""

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(nil, errors.New("not found"))

	// Act
	result, err := s.service.CancelOrderByUuid(ctx, orderUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.ErrOrderNotFound, err)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *CancelOrderTestSuite) TestCancelOrderByUuid_WithTransactionUUID() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "550e8400-e29b-41d4-a716-446655440001"
	transactionUUID := "550e8400-e29b-41d4-a716-446655440003"

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   repoModel.PaymentMethodSBP,
		Status:          repoModel.StatusPendingPayment,
	}

	expectedModelOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   model.PaymentMethodSBP,
		Status:          model.StatusCanceled,
	}

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   repoModel.PaymentMethodSBP,
		Status:          repoModel.StatusCanceled,
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(repoOrder, nil)
	s.orderRepository.On("UpdateOrder", ctx, expectedRepoOrder).Return(nil)

	// Act
	result, err := s.service.CancelOrderByUuid(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedModelOrder, result)
}
