package producer

import (
	"context"

	"google.golang.org/protobuf/proto"

	orderKafka "github.com/space-wanderer/microservices/order/internal/converter/kafka"
	"github.com/space-wanderer/microservices/order/internal/model"
	platformKafka "github.com/space-wanderer/microservices/platform/pkg/kafka"
	events_v1 "github.com/space-wanderer/microservices/shared/pkg/proto/events/v1"
)

type orderProducer struct {
	producer platformKafka.Producer
}

func NewOrderProducer(producer platformKafka.Producer) orderKafka.OrderPaidProducer {
	return &orderProducer{
		producer: producer,
	}
}

func (p *orderProducer) ProduceOrderPaidEvent(ctx context.Context, event model.OrderPaidEvent) error {
	pbEvent := &events_v1.OrderPaidEvent{
		EventUuid:       event.EventUUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID,
	}

	data, err := proto.Marshal(pbEvent)
	if err != nil {
		return err
	}

	return p.producer.Send(ctx, []byte(event.OrderUUID), data)
}
