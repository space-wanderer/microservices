package order

import (
	"context"

	"github.com/space-wanderer/microservices/order/internal/converter"
	"github.com/space-wanderer/microservices/order/internal/model"
)

func (s *service) GetOrderByUuid(ctx context.Context, uuid string) (model.Order, error) {
	orderModel, err := s.orderRepository.GetOrderByUuid(ctx, uuid)
	if err != nil {
		return model.Order{}, model.ErrOrderNotFound
	}

	order := converter.ConvertRepoOrderToModelOrder(orderModel)
	if order == nil {
		return model.Order{}, model.ErrOrderNotFound
	}

	return *order, nil
}
