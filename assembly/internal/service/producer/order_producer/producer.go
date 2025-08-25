package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/space-wanderer/microservices/assembly/internal/model"
	"github.com/space-wanderer/microservices/platform/pkg/kafka"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
	eventsV1 "github.com/space-wanderer/microservices/shared/pkg/proto/events/v1"
)

type service struct {
	assemblyRecodedProducer kafka.Producer
}

func NewService(assemblyRecodedProducer kafka.Producer) *service {
	return &service{
		assemblyRecodedProducer: assemblyRecodedProducer,
	}
}

func (s *service) ProduceShipAssembledEvent(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := &eventsV1.ShipAssembledEvent{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal order paid event", zap.Error(err))
		return err
	}

	err = s.assemblyRecodedProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to send order paid event", zap.Error(err))
		return err
	}

	return nil
}
