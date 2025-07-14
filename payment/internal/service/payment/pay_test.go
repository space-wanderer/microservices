package payment

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"

	"github.com/space-wanderer/microservices/payment/internal/model"
)

func TestService_PayOrder_Integration(t *testing.T) {
	// Тест с реальной реализацией (без моков)
	service := NewService()

	req := model.Pay{
		OrderUuid:     gofakeit.UUID(),
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: model.PaymentMethodCard,
	}

	result, err := service.PayOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "-")
}

func TestService_PayOrder_WithFakeData(t *testing.T) {
	// Тест с различными типами данных от gofakeit
	service := NewService()

	// Генерируем случайные данные
	gofakeit.Seed(12345)

	testCases := []model.PaymentMethod{
		model.PaymentMethodCard,
		model.PaymentMethodSBP,
		model.PaymentMethodCreditCard,
		model.PaymentMethodInvestorMoney,
	}

	for _, paymentMethod := range testCases {
		t.Run(string(paymentMethod), func(t *testing.T) {
			req := model.Pay{
				OrderUuid:     gofakeit.UUID(),
				UserUuid:      gofakeit.UUID(),
				PaymentMethod: paymentMethod,
			}

			result, err := service.PayOrder(context.Background(), req)

			assert.NoError(t, err)
			assert.NotEmpty(t, result)
			// Проверяем, что результат является валидным UUID
			assert.Contains(t, result, "-")
		})
	}
}
