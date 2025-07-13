package payment

import (
	"github.com/space-wanderer/microservices/payment/internal/service"
)

type Service struct {
	paymentService service.PaymentService
}

func NewService(paymentService service.PaymentService) *Service {
	return &Service{paymentService: paymentService}
}
