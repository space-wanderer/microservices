package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/space-wanderer/microservices/order/internal/converter/kafka"
	"github.com/space-wanderer/microservices/order/internal/model"
	events_v1 "github.com/space-wanderer/microservices/shared/pkg/proto/events/v1"
)

type shipAssembledDecoder struct{}

func NewShipAssembledDecoder() kafka.ShipAssembledDecoder {
	return &shipAssembledDecoder{}
}

func (d *shipAssembledDecoder) Decode(data []byte) model.ShipAssembledEvent {
	var pbEvent events_v1.ShipAssembledEvent
	if err := proto.Unmarshal(data, &pbEvent); err != nil {
		return model.ShipAssembledEvent{}
	}

	return model.ShipAssembledEvent{
		EventUUID:    pbEvent.EventUuid,
		OrderUUID:    pbEvent.OrderUuid,
		UserUUID:     pbEvent.UserUuid,
		BuildTimeSec: int(pbEvent.BuildTimeSec),
	}
}
