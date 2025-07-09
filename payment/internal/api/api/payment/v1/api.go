package v1

import (
	"github.com/space-wanderer/microservices/payment/internal/service"

	paymentV1 "github.com/space-wanderer/microservices/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{paymentService: paymentService}
}
