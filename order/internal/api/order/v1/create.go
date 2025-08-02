package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/space-wanderer/microservices/order/internal/converter"
	orderV1 "github.com/space-wanderer/microservices/shared/pkg/api/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	// Конвертируем запрос в модель сервиса
	order := converter.ConvertCreateOrderRequestToModelOrder(req)

	// Создаем заказ через сервис
	createdOrder, err := a.orderService.CreateOrder(ctx, *order)
	if err != nil {
		return nil, err
	}

	// Конвертируем ответ
	orderUUID, err := uuid.Parse(createdOrder.OrderUUID)
	if err != nil {
		return nil, err
	}
	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: createdOrder.TotalPrice,
	}, nil
}
