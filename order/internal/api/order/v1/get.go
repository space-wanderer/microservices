package v1

import (
	"context"

	"github.com/space-wanderer/microservices/order/internal/converter"
	orderV1 "github.com/space-wanderer/microservices/shared/pkg/api/order/v1"
)

func (a *api) GetOrderByUuid(ctx context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	// Получаем заказ через сервис
	order, err := a.orderService.GetOrderByUuid(ctx, params.OrderUUID.String())
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Заказ не найден",
		}, nil
	}

	// Конвертируем в DTO
	orderDto := converter.ConvertModelOrderToOrderDto(&order)

	return &orderV1.GetOrderResponse{
		Order:   orderDto,
		Message: orderV1.NewOptString("Заказ успешно получен"),
	}, nil
}
