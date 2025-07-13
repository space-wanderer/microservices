package order

import (
	"context"
	"fmt"

	"github.com/space-wanderer/microservices/order/internal/repository/model"
)

func (r *repository) GetOrderByUuid(ctx context.Context, uuid string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[uuid]
	if !exists {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}
