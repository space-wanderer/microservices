package converter

import (
	"github.com/space-wanderer/microservices/payment/internal/model"

	paymentV1 "github.com/space-wanderer/microservices/shared/pkg/proto/payment/v1"
)

// convertFromGRPC конвертирует gRPC запрос во внутреннюю модель
func ConvertFromGRPC(req *paymentV1.PayOrderRequest) model.Pay {
	return model.Pay{
		OrderUuid:     req.OrderUuid,
		UserUuid:      req.UserUuid,
		PaymentMethod: convertPaymentMethod(req.PaymentMethod),
	}
}

// convertPaymentMethod конвертирует gRPC PaymentMethod во внутренний
func convertPaymentMethod(method paymentV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}
