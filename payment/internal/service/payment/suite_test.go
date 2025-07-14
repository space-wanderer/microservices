package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	cnx     context.Context
	service *Service
}

func (s *ServiceSuite) SetupTest() {
	s.cnx = context.Background()
	s.service = NewService()
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
