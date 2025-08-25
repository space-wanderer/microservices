package kafka

import (
	"github.com/space-wanderer/microservices/notification/internal/model"
)

type OrderPaidDecoder interface {
	Decode(data []byte) model.OrderPaidEvent
}

type ShipAssembledDecoder interface {
	Decode(data []byte) model.ShipAssembledEvent
}
