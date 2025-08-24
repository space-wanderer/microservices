package order

import (
	"context"
	"fmt"

	"github.com/space-wanderer/microservices/order/internal/client/grpc"
	"github.com/space-wanderer/microservices/order/internal/converter"
	kafkaConverter "github.com/space-wanderer/microservices/order/internal/converter/kafka"
	"github.com/space-wanderer/microservices/order/internal/model"
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

func (s *service) UpdateOrderStatus(ctx context.Context, orderUUID string, status model.Status) error {
	// Получаем заказ из репозитория
	repoOrder, err := s.orderRepository.GetOrderByUuid(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Конвертируем в модель сервиса
	order := converter.ConvertRepoOrderToModelOrder(repoOrder)

	// Обновляем статус
	order.Status = status

	// Конвертируем обратно в модель репозитория и сохраняем
	repoOrder = converter.ConvertModelOrderToRepoOrder(order)

	err = s.orderRepository.UpdateOrder(ctx, repoOrder)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}
