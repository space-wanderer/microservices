package order

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/space-wanderer/microservices/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite
	orderRepository *mocks.OrderRepository
}

func (s *ServiceSuite) SetupTest() {
	s.orderRepository = mocks.NewOrderRepository(s.T())
}

func (s *ServiceSuite) TearDownTest() {
	s.orderRepository.AssertExpectations(s.T())
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
