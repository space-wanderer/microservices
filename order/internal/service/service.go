package service

import (
	"context"

	"github.com/space-wanderer/microservices/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req model.Order) (model.Order, error)
	GetOrderByUuid(ctx context.Context, uuid string) (model.Order, error)
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (model.Order, error)
	CancelOrderByUuid(ctx context.Context, uuid string) (model.Order, error)
	UpdateOrderStatus(ctx context.Context, orderUUID string, status model.Status) error
}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
