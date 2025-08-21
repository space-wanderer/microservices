package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/space-wanderer/microservices/notification/internal/model"
	eventsV1 "github.com/space-wanderer/microservices/shared/pkg/proto/events/v1"
)

type orderPaidDecoder struct{}

func NewOrderPaidDecoder() *orderPaidDecoder {
	return &orderPaidDecoder{}
}

func (d *orderPaidDecoder) Decode(data []byte) model.OrderPaidEvent {
	var pb eventsV1.OrderPaidEvent
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}
	}

	return model.OrderPaidEvent{
		EventUUID:       pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		PaymentMethod:   pb.PaymentMethod,
		TransactionUUID: pb.TransactionUuid,
	}
}
