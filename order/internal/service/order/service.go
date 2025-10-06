package order

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

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

	// МЕТРИКИ: счетчики для мониторинга бизнес-метрик
	ordersTotal        metric.Int64Counter   // количество созданных заказов
	ordersRevenueTotal metric.Float64Counter // суммарная выручка от заказов
}

func NewOrderService(orderRepository repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient, orderPaidProducer kafkaConverter.OrderPaidProducer) *service {
	// ИНИЦИАЛИЗАЦИЯ МЕТРИК: создаем счетчики для мониторинга
	meter := otel.Meter("order-service")

	ordersTotal, err := meter.Int64Counter(
		"orders_total",
		metric.WithDescription("Total number of orders created"),
	)
	if err != nil {
		// Логируем ошибку, но продолжаем работу
		// В production можно использовать fallback метрику
		_ = err // Игнорируем ошибку для graceful degradation
	}

	ordersRevenueTotal, err := meter.Float64Counter(
		"orders_revenue_total",
		metric.WithDescription("Total revenue from orders"),
	)
	if err != nil {
		// Логируем ошибку, но продолжаем работу
		// В production можно использовать fallback метрику
		_ = err // Игнорируем ошибку для graceful degradation
	}

	return &service{
		orderRepository:    orderRepository,
		inventoryClient:    inventoryClient,
		paymentClient:      paymentClient,
		orderPaidProducer:  orderPaidProducer,
		ordersTotal:        ordersTotal,
		ordersRevenueTotal: ordersRevenueTotal,
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
