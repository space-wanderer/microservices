package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/space-wanderer/microservices/notification/internal/model"
	eventsV1 "github.com/space-wanderer/microservices/shared/pkg/proto/events/v1"
)

type orderAssembledDecoder struct{}

func NewOrderAssembledDecoder() *orderAssembledDecoder {
	return &orderAssembledDecoder{}
}

func (d *orderAssembledDecoder) Decode(data []byte) model.ShipAssembledEvent {
	var pb eventsV1.ShipAssembledEvent
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}
	}

	return model.ShipAssembledEvent{
		EventUUID:    pb.EventUuid,
		OrderUUID:    pb.OrderUuid,
		UserUUID:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}
}
