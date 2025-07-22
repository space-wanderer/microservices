package repository

import (
	"context"

	"github.com/space-wanderer/microservices/order/internal/repository/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) (string, error)
	GetOrderByUuid(ctx context.Context, uuid string) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
}
