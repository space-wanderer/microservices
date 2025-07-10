package v1

import (
	"context"

	"github.com/space-wanderer/microservices/payment/internal/converter"
	"github.com/space-wanderer/microservices/payment/internal/model"
	paymentV1 "github.com/space-wanderer/microservices/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	// Конвертируем gRPC запрос во внутреннюю модель
	payment := converter.ConvertFromGRPC(req)

	transactionUUID, err := a.paymentService.PayOrder(ctx, payment)
	if err != nil {
		return nil, model.ErrPayment
	}

	// Возвращаем gRPC ответ
	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
