package part

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/space-wanderer/microservices/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite
	inventoryRepository *mocks.InventoryRepository
	service             *Service
}

func (s *ServiceSuite) SetupTest() {
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
