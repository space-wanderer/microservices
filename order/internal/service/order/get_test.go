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

type GetOrderTestSuite struct {
	suite.Suite
	orderRepository *repoMocks.OrderRepository
	inventoryClient *mocks.InventoryClient
	paymentClient   *mocks.PaymentClient
	service         *service
}

func (s *GetOrderTestSuite) SetupTest() {
	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = mocks.NewInventoryClient(s.T())
	s.paymentClient = mocks.NewPaymentClient(s.T())
	s.service = NewOrderService(s.orderRepository, s.inventoryClient, s.paymentClient)
}

func (s *GetOrderTestSuite) TearDownTest() {
	s.orderRepository.AssertExpectations(s.T())
	s.inventoryClient.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

func TestGetOrderTestSuite(t *testing.T) {
	suite.Run(t, new(GetOrderTestSuite))
}

func (s *GetOrderTestSuite) TestGetOrderByUuid_Success() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        "550e8400-e29b-41d4-a716-446655440001",
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      150.5,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCard,
		Status:          repoModel.StatusPendingPayment,
	}

	expectedModelOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        expectedRepoOrder.UserUUID,
		PartUuids:       expectedRepoOrder.PartUuids,
		TotalPrice:      expectedRepoOrder.TotalPrice,
		TransactionUUID: nil,
		PaymentMethod:   model.PaymentMethodCard,
		Status:          model.StatusPendingPayment,
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(expectedRepoOrder, nil)

	// Act
	result, err := s.service.GetOrderByUuid(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedModelOrder, result)
}

func (s *GetOrderTestSuite) TestGetOrderByUuid_RepositoryError() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	expectedError := errors.New("database error")

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(nil, expectedError)

	// Act
	result, err := s.service.GetOrderByUuid(ctx, orderUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *GetOrderTestSuite) TestGetOrderByUuid_EmptyUUID() {
	// Arrange
	ctx := context.Background()
	orderUUID := ""

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(nil, errors.New("invalid uuid"))

	// Act
	result, err := s.service.GetOrderByUuid(ctx, orderUUID)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), model.Order{}, result)
}

func (s *GetOrderTestSuite) TestGetOrderByUuid_WithTransactionUUID() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"
	transactionUUID := "550e8400-e29b-41d4-a716-446655440003"

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        "550e8400-e29b-41d4-a716-446655440001",
		PartUuids:       []string{"550e8400-e29b-41d4-a716-446655440002"},
		TotalPrice:      250.75,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   repoModel.PaymentMethodSBP,
		Status:          repoModel.StatusPaid,
	}

	expectedModelOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        expectedRepoOrder.UserUUID,
		PartUuids:       expectedRepoOrder.PartUuids,
		TotalPrice:      expectedRepoOrder.TotalPrice,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   model.PaymentMethodSBP,
		Status:          model.StatusPaid,
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(expectedRepoOrder, nil)

	// Act
	result, err := s.service.GetOrderByUuid(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedModelOrder, result)
}

func (s *GetOrderTestSuite) TestGetOrderByUuid_MultipleParts() {
	// Arrange
	ctx := context.Background()
	orderUUID := "550e8400-e29b-41d4-a716-446655440000"

	partUUIDs := []string{
		"550e8400-e29b-41d4-a716-446655440002",
		"550e8400-e29b-41d4-a716-446655440003",
		"550e8400-e29b-41d4-a716-446655440004",
	}

	expectedRepoOrder := &repoModel.Order{
		OrderUUID:       orderUUID,
		UserUUID:        "550e8400-e29b-41d4-a716-446655440001",
		PartUuids:       partUUIDs,
		TotalPrice:      450.25,
		TransactionUUID: nil,
		PaymentMethod:   repoModel.PaymentMethodCreditCard,
		Status:          repoModel.StatusPendingPayment,
	}

	expectedModelOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        expectedRepoOrder.UserUUID,
		PartUuids:       expectedRepoOrder.PartUuids,
		TotalPrice:      expectedRepoOrder.TotalPrice,
		TransactionUUID: nil,
		PaymentMethod:   model.PaymentMethodCreditCard,
		Status:          model.StatusPendingPayment,
	}

	s.orderRepository.On("GetOrderByUuid", ctx, orderUUID).Return(expectedRepoOrder, nil)

	// Act
	result, err := s.service.GetOrderByUuid(ctx, orderUUID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedModelOrder, result)
	assert.Len(s.T(), result.PartUuids, 3)
}
