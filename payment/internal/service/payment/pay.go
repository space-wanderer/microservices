package payment

import (
	"context"

	"github.com/space-wanderer/microservices/payment/internal/model"
)

func (s *Service) PayOrder(ctx context.Context, req model.Pay) (string, error) {
	transactionUUID, err := s.paymentService.PayOrder(ctx, req)
	if err != nil {
		return "", model.ErrPayment
	}

	return transactionUUID, nil
}
