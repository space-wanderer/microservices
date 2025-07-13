package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/space-wanderer/microservices/order/internal/converter"
	"github.com/space-wanderer/microservices/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (model.Order, error) {
	repoOrder, err := s.orderRepository.GetOrderByUuid(ctx, orderUUID)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to get order: %w", err)
	}

	// Конвертируем в модель сервиса
	order := converter.ConvertRepoOrderToModelOrder(repoOrder)

	// Проверяем статус заказа
	if order.Status != model.StatusPendingPayment {
		return model.Order{}, errors.New("order is not in pending payment status")
	}

	// Обрабатываем платеж через PaymentService
	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID, userUUID, string(paymentMethod))
	if err != nil {
		return model.Order{}, fmt.Errorf("payment processing failed: %w", err)
	}

	// Обновляем заказ после успешного платежа
	order.Status = model.StatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = paymentMethod

	// Конвертируем обратно в модель репозитория и сохраняем
	repoOrder = converter.ConvertModelOrderToRepoOrder(order)
	err = s.orderRepository.UpdateOrder(ctx, repoOrder)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to update order: %w", err)
	}

	return *order, nil
}
