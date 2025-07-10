package service

import (
	"context"

	"github.com/space-wanderer/microservices/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, req model.Pay) (string, error)
}
