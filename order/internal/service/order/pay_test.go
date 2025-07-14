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

type PayOrderTestSuite struct {
	suite.Suite
	orderRepository *repoMocks.OrderRepository
	inventoryClient *mocks.InventoryClient
	paymentClient   *mocks.PaymentClient
	service         *service
}

func (s *PayOrderTestSuite) SetupTest() {
	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = mocks.NewInventoryClient(s.T())
	s.paymentClient = mocks.NewPaymentClient(s.T())
	s.service = NewOrderService(s.orderRepository, s.inventoryClient, s.paymentClient)
}

func (s *PayOrderTestSuite) TearDownTest() {
	s.orderRepository.AssertExpectations(s.T())
	s.inventoryClient.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

func TestPayOrderTestSuite(t *testing.T) {
	suite.Run(t, new(PayOrderTestSuite))
}

func (s *PayOrderTestSuite) TestPayOrder_Success() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "550e8400-e29b-41d4-a716-446655440001"
	transactionUUID := "550e8400-e29b-41d4-a716-446655440003"
	paymentMethod := model.PaymentMethodCard

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
		TransactionUUID: &transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          model.StatusPaid,
	}

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPaid,
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(repoOrder, nil)
	s.paymentClient.On("PayOrder", ctx, orderUUID, userUUID, string(paymentMethod)).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", ctx, expectedRepoOrder).Return(nil)

	// Act
	result, err := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedModelOrder, result)
}

func (s *PayOrderTestSuite) TestPayOrder_GetOrderError() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "550e8400-e29b-41d4-a716-446655440001"
	paymentMethod := model.PaymentMethodCard
	expectedError := errors.New("database error")

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(nil, expectedError)

	// Act
	result, err := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to get order")
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *PayOrderTestSuite) TestPayOrder_OrderNotPendingPayment() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "550e8400-e29b-41d4-a716-446655440001"
	paymentMethod := model.PaymentMethodCard

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPaid, // Уже оплачен
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(repoOrder, nil)

	// Act
	result, err := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), "order is not in pending payment status", err.Error())
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *PayOrderTestSuite) TestPayOrder_PaymentClientError() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "550e8400-e29b-41d4-a716-446655440001"
	paymentMethod := model.PaymentMethodCard
	expectedError := errors.New("payment service error")

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPendingPayment,
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(repoOrder, nil)
	s.paymentClient.On("PayOrder", ctx, orderUUID, userUUID, string(paymentMethod)).Return("", expectedError)

	// Act
	result, err := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "payment processing failed")
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *PayOrderTestSuite) TestPayOrder_UpdateOrderError() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "550e8400-e29b-41d4-a716-446655440001"
	transactionUUID := "550e8400-e29b-41d4-a716-446655440003"
	paymentMethod := model.PaymentMethodCard
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
		TransactionUUID: &transactionUUID,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPaid,
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(repoOrder, nil)
	s.paymentClient.On("PayOrder", ctx, orderUUID, userUUID, string(paymentMethod)).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", ctx, expectedRepoOrder).Return(expectedError)

	// Act
	result, err := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to update order")
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *PayOrderTestSuite) TestPayOrder_DifferentPaymentMethod() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "550e8400-e29b-41d4-a716-446655440001"
	transactionUUID := "550e8400-e29b-41d4-a716-446655440003"
	paymentMethod := model.PaymentMethodSBP

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
		TransactionUUID: &transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          model.StatusPaid,
	}

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   repoModel.PaymentMethodSBP,
		Status:          repoModel.StatusPaid,
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(repoOrder, nil)
	s.paymentClient.On("PayOrder", ctx, orderUUID, userUUID, string(paymentMethod)).Return(transactionUUID, nil)
	s.orderRepository.On("UpdateOrder", ctx, expectedRepoOrder).Return(nil)

	// Act
	result, err := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedModelOrder, result)
}

func (s *PayOrderTestSuite) TestPayOrder_CanceledOrder() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	userUUID := "550e8400-e29b-41d4-a716-446655440001"
	paymentMethod := model.PaymentMethodCard

	repoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusCanceled, // Отменен
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(repoOrder, nil)

	// Act
	result, err := s.service.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), "order is not in pending payment status", err.Error())
	assert.Equal(s.T(), model.Order{}, result)
}
