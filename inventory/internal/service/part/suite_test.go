package part

import (
	"context"
	"testing"

	"github.com/space-wanderer/microservices/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	cnx                 context.Context
	inventoryRepository *mocks.InventoryRepository
	service             *service
}

func (s *ServiceSuite) SetupTest() {
	s.cnx = context.Background()
	s.inventoryRepository = mocks.NewInventoryRepository(s.T())
	s.service = NewService(
		s.inventoryRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
	s.inventoryRepository.AssertExpectations(s.T())
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
