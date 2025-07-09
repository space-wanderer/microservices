package service

import (
	"context"
)

type PaymentService interface {
	PayOrder(ctx context.Context, uuid string) error
}
