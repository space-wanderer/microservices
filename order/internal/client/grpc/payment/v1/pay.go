package v1

import (
	"context"
	"time"

	generatedPaymentV1 "github.com/space-wanderer/microservices/shared/pkg/proto/payment/v1"
)

// PayOrder обрабатывает платеж через PaymentService
func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
	var grpcPaymentMethod generatedPaymentV1.PaymentMethod
	switch paymentMethod {
	case "CARD":
		grpcPaymentMethod = generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case "SBP":
		grpcPaymentMethod = generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case "CREDIT_CARD":
		grpcPaymentMethod = generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case "INVESTOR_MONEY":
		grpcPaymentMethod = generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		grpcPaymentMethod = generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	req := &generatedPaymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: grpcPaymentMethod,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.generatedClient.PayOrder(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.TransactionUuid, nil
}
