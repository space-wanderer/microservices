package order

import (
	"context"
	"fmt"

	"github.com/space-wanderer/microservices/order/internal/converter"
	"github.com/space-wanderer/microservices/order/internal/model"
)

func (s *service) CancelOrderByUuid(ctx context.Context, orderUUID string) (model.Order, error) {
	repoOrder, err := s.orderRepository.GetOrderByUuid(ctx, orderUUID)
	if err != nil {
		return model.Order{}, model.ErrOrderNotFound
	}

	// Конвертируем в модель сервиса
	order := converter.ConvertRepoOrderToModelOrder(repoOrder)

	if order.Status == model.StatusPaid {
		return model.Order{}, model.ErrOrderCannotBeCancelled
	}

	order.Status = model.StatusCanceled

	// Конвертируем обратно в модель репозитория и сохраняем
	repoOrder = converter.ConvertModelOrderToRepoOrder(order)
	err = s.orderRepository.UpdateOrder(ctx, repoOrder)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to update order: %w", err)
	}

	return *order, nil
}
