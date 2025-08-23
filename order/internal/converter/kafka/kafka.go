package kafka

import (
	"context"

	"github.com/space-wanderer/microservices/order/internal/model"
)

// OrderPaidProducer интерфейс для отправки OrderPaidEvent
type OrderPaidProducer interface {
	ProduceOrderPaidEvent(ctx context.Context, event model.OrderPaidEvent) error
}
