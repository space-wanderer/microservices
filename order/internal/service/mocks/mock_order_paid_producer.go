package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/space-wanderer/microservices/order/internal/model"
)

// MockOrderPaidProducer is a mock of OrderPaidProducer interface.
type MockOrderPaidProducer struct {
	mock.Mock
}

// NewMockOrderPaidProducer creates a new mock instance.
func NewMockOrderPaidProducer(t mock.TestingT) *MockOrderPaidProducer {
	mock := &MockOrderPaidProducer{}
	mock.Test(t)
	return mock
}

// ProduceOrderPaidEvent mocks base method.
func (m *MockOrderPaidProducer) ProduceOrderPaidEvent(ctx context.Context, event model.OrderPaidEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
