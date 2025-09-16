package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/space-wanderer/microservices/payment/internal/model"
)

func (s *Service) PayOrder(ctx context.Context, req model.Pay) (string, error) {
	// ТРЕЙСИНГ: продолжаем трейс из OrderService
	tracer := otel.Tracer("payment-service")
	ctx, span := tracer.Start(ctx, "ProcessPayment")
	defer span.End()

	// Добавляем атрибуты к span
	span.SetAttributes(
		attribute.String("payment.order_id", req.OrderUuid),
		attribute.String("payment.user_id", req.UserUuid),
		attribute.String("payment.method", string(req.PaymentMethod)),
	)

	transactionUUID := uuid.New().String()

	// Имитируем обработку платежа
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	// Добавляем результат в span
	span.SetAttributes(
		attribute.String("payment.transaction_id", transactionUUID),
		attribute.String("payment.result", "success"),
		attribute.String("payment.status", "completed"),
	)

	return transactionUUID, nil
}
