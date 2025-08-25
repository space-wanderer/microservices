package decoder

import (
	"google.golang.org/protobuf/proto"

	"github.com/space-wanderer/microservices/assembly/internal/model"
	eventsV1 "github.com/space-wanderer/microservices/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) model.OrderPaidEvent {
	var pb eventsV1.OrderPaidEvent
	if err := proto.Unmarshal(data, &pb); err != nil {
		// В случае ошибки возвращаем пустую структуру
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
