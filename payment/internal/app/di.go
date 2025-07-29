package app

import (
	"context"

	"github.com/space-wanderer/microservices/payment/internal/service"
	"github.com/space-wanderer/microservices/payment/internal/service/payment"
)

type diContainer struct {
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = payment.NewService()
	}
	return d.paymentService
}
