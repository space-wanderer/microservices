package order

import (
	"context"

	"github.com/space-wanderer/microservices/order/internal/repository/model"
)

func (r *repository) UpdateOrder(ctx context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.OrderUUID] = order
	return nil
}
