package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/space-wanderer/microservices/order/internal/converter"
	"github.com/space-wanderer/microservices/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, req model.Order) (model.Order, error) {
	repoOrder := converter.ConvertModelOrderToRepoOrder(&req)
	orderUUID, err := s.orderRepository.CreateOrder(ctx, repoOrder)
	if err != nil {
		return model.Order{}, err
	}

	// Конвертируем []string в []uuid.UUID через converter
	partUUIDs, err := converter.ConvertStringSliceToUUIDSlice(req.PartUuids)
	if err != nil {
		return model.Order{}, fmt.Errorf("ошибка при конвертации UUID: %w", err)
	}

	totalPrice, err := s.calculateOrderPrice(ctx, partUUIDs)
	if err != nil {
		return model.Order{}, fmt.Errorf("ошибка при получении информации о деталях: %w", err)
	}

	return model.Order{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

func (s *service) calculateOrderPrice(ctx context.Context, partUuids []uuid.UUID) (float32, error) {
	var totalPrice float32

	for _, partUUID := range partUuids {
		filter := model.PartsFilter{
			Uuids: []string{partUUID.String()},
		}

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		parts, err := s.inventoryClient.ListParts(ctx, filter)
		if err != nil {
			return 0, err
		}

		if len(parts) == 0 {
			return 0, fmt.Errorf("part not found: %s", partUUID.String())
		}

		totalPrice += float32(parts[0].Price)
	}

	return totalPrice, nil
}
