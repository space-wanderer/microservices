package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/space-wanderer/microservices/order/internal/model"
	orderV1 "github.com/space-wanderer/microservices/shared/pkg/api/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order, err := a.orderService.PayOrder(ctx, params.OrderUUID.String(), "", model.PaymentMethod(req.PaymentMethod))
	if err != nil {
		return &orderV1.BadGatewayError{
			Error:   "PAYMENT_ERROR",
			Message: err.Error(),
		}, nil
	}

	if order.TransactionUUID == nil {
		return &orderV1.BadGatewayError{
			Error:   "PAYMENT_ERROR",
			Message: "Transaction UUID not found",
		}, nil
	}

	transactionUUID, err := uuid.Parse(*order.TransactionUUID)
	if err != nil {
		return &orderV1.BadGatewayError{
			Error:   "PAYMENT_ERROR",
			Message: "Invalid transaction UUID format",
		}, nil
	}
	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
