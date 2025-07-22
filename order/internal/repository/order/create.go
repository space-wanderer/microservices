package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/space-wanderer/microservices/order/internal/repository/model"
)

func (r *repository) CreateOrder(_ context.Context, req *model.Order) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	orderUUID := uuid.New().String()

	order := &model.Order{
		OrderUUID:  orderUUID,
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: req.TotalPrice,
	}

	r.orders[orderUUID] = order
	return order.OrderUUID, nil
}
