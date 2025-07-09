package payment

import def "github.com/space-wanderer/microservices/payment/internal/service"

var _ def.PaymentService = (*Service)(nil)

type service struct {
	paymentService service.PaymentService
}

func NewService(paymentService service.PaymentService) *service {
	return &service{paymentService: paymentService}
}
