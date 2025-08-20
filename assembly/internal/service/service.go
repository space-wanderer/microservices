package service

import (
	"context"

	"github.com/space-wanderer/microservices/assembly/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ProducerService interface {
	ProduceShipAssembledEvent(ctx context.Context, event model.ShipAssembledEvent) error
}
