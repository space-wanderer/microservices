package service

import (
	"context"

	"github.com/space-wanderer/microservices/notification/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, uuid string, event model.OrderPaidEvent) error
	SendShipAssembledNotification(ctx context.Context, uuid string, event model.ShipAssembledEvent) error
}
