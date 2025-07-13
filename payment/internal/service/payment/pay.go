package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/space-wanderer/microservices/payment/internal/model"
)

func (s *Service) PayOrder(ctx context.Context, req model.Pay) (string, error) {
	transactionUUID := uuid.New().String()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	return transactionUUID, nil
}
