package v1

import (
	"context"
	"time"

	generatedPaymentV1 "github.com/space-wanderer/microservices/shared/pkg/proto/payment/v1"
)

// PayOrder обрабатывает платеж через PaymentService
func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {

	req := &generatedPaymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: generatedPaymentV1.PaymentMethod(generatedPaymentV1.PaymentMethod_value[paymentMethod]),
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.generatedClient.PayOrder(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.TransactionUuid, nil
}
