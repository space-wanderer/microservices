package kafka

import (
	"context"

	"github.com/space-wanderer/microservices/order/internal/model"
)

// OrderPaidProducer интерфейс для отправки OrderPaidEvent
type OrderPaidProducer interface {
	ProduceOrderPaidEvent(ctx context.Context, event model.OrderPaidEvent) error
}

// ShipAssembledDecoder интерфейс для декодирования ShipAssembledEvent
type ShipAssembledDecoder interface {
	Decode(data []byte) model.ShipAssembledEvent
}
