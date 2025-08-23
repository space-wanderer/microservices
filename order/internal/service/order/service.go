package order

import (
	"github.com/space-wanderer/microservices/order/internal/client/grpc"
	kafkaConverter "github.com/space-wanderer/microservices/order/internal/converter/kafka"
	"github.com/space-wanderer/microservices/order/internal/repository"
)

type service struct {
	orderRepository   repository.OrderRepository
	inventoryClient   grpc.InventoryClient
	paymentClient     grpc.PaymentClient
	orderPaidProducer kafkaConverter.OrderPaidProducer
}

func NewOrderService(orderRepository repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient, orderPaidProducer kafkaConverter.OrderPaidProducer) *service {
	return &service{
		orderRepository:   orderRepository,
		inventoryClient:   inventoryClient,
		paymentClient:     paymentClient,
		orderPaidProducer: orderPaidProducer,
	}
}
