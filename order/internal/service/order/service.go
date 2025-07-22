package order

import (
	"github.com/space-wanderer/microservices/order/internal/client/grpc"
	"github.com/space-wanderer/microservices/order/internal/repository"
)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewOrderService(orderRepository repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
