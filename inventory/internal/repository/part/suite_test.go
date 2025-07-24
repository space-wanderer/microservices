package part

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/space-wanderer/microservices/inventory/internal/repository/mocks"
)

type RepositorySuite struct {
	suite.Suite
	mockRepository *mocks.InventoryRepository
}

func (s *RepositorySuite) SetupTest() {
	s.mockRepository = mocks.NewInventoryRepository(s.T())
}

func (s *RepositorySuite) TearDownTest() {
	s.mockRepository.AssertExpectations(s.T())
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
