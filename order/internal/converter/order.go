package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/space-wanderer/microservices/order/internal/model"
	repoModel "github.com/space-wanderer/microservices/order/internal/repository/model"
	order_v1 "github.com/space-wanderer/microservices/shared/pkg/api/order/v1"
)

// ConvertRepoOrderToModelOrder конвертирует Order из repository model в service model
func ConvertRepoOrderToModelOrder(repoOrder *repoModel.Order) *model.Order {
	if repoOrder == nil {
		return nil
	}

	return &model.Order{
		OrderUUID:       repoOrder.OrderUUID,
		UserUUID:        repoOrder.UserUUID,
		PartUuids:       repoOrder.PartUuids,
		TotalPrice:      repoOrder.TotalPrice,
		TransactionUUID: repoOrder.TransactionUUID,
		PaymentMethod:   convertRepoPaymentMethodToModelPaymentMethod(repoOrder.PaymentMethod),
		Status:          convertRepoStatusToModelStatus(repoOrder.Status),
	}
}

// ConvertModelOrderToRepoOrder конвертирует Order из service model в repository model
func ConvertModelOrderToRepoOrder(modelOrder *model.Order) *repoModel.Order {
	if modelOrder == nil {
		return nil
	}

	return &repoModel.Order{
		OrderUUID:       modelOrder.OrderUUID,
		UserUUID:        modelOrder.UserUUID,
		PartUuids:       modelOrder.PartUuids,
		TotalPrice:      modelOrder.TotalPrice,
		TransactionUUID: modelOrder.TransactionUUID,
		PaymentMethod:   convertModelPaymentMethodToRepoPaymentMethod(modelOrder.PaymentMethod),
		Status:          convertModelStatusToRepoStatus(modelOrder.Status),
	}
}

// ConvertModelOrderToOrderDto конвертирует Order из service model в API DTO
func ConvertModelOrderToOrderDto(modelOrder *model.Order) order_v1.OrderDto {
	orderDto := order_v1.OrderDto{
		OrderUUID:  uuid.MustParse(modelOrder.OrderUUID),
		UserUUID:   uuid.MustParse(modelOrder.UserUUID),
		PartUuids:  convertStringSliceToUUIDSlice(modelOrder.PartUuids),
		TotalPrice: modelOrder.TotalPrice,
		Status:     convertModelStatusToOrderStatus(modelOrder.Status),
	}

	if modelOrder.TransactionUUID != nil {
		transactionUUID := uuid.MustParse(*modelOrder.TransactionUUID)
		orderDto.TransactionUUID = order_v1.NewOptUUID(transactionUUID)
	}

	if modelOrder.PaymentMethod != "" {
		orderDto.PaymentMethod = convertModelPaymentMethodToOrderPaymentMethod(modelOrder.PaymentMethod)
	}

	return orderDto
}

// ConvertCreateOrderRequestToModelOrder конвертирует CreateOrderRequest в service model
func ConvertCreateOrderRequestToModelOrder(req *order_v1.CreateOrderRequest) *model.Order {
	return &model.Order{
		UserUUID:  req.UserUUID.String(),
		PartUuids: convertUUIDSliceToStringSlice(req.PartUuids),
		Status:    model.StatusPendingPayment,
	}
}

// convertRepoPaymentMethodToModelPaymentMethod конвертирует PaymentMethod из repository в service model
func convertRepoPaymentMethodToModelPaymentMethod(repoMethod repoModel.PaymentMethod) model.PaymentMethod {
	switch repoMethod {
	case repoModel.PaymentMethodCard:
		return model.PaymentMethodCard
	case repoModel.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case repoModel.PaymentMethodCreditCard:
		return model.PaymentMethodCreditCard
	case repoModel.PaymentMethodInvestorMoney:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}

// convertModelPaymentMethodToRepoPaymentMethod конвертирует PaymentMethod из service в repository model
func convertModelPaymentMethodToRepoPaymentMethod(modelMethod model.PaymentMethod) repoModel.PaymentMethod {
	switch modelMethod {
	case model.PaymentMethodCard:
		return repoModel.PaymentMethodCard
	case model.PaymentMethodSBP:
		return repoModel.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return repoModel.PaymentMethodCreditCard
	case model.PaymentMethodInvestorMoney:
		return repoModel.PaymentMethodInvestorMoney
	default:
		return repoModel.PaymentMethodUnknown
	}
}

// convertRepoStatusToModelStatus конвертирует Status из repository в service model
func convertRepoStatusToModelStatus(repoStatus repoModel.Status) model.Status {
	switch repoStatus {
	case repoModel.StatusPendingPayment:
		return model.StatusPendingPayment
	case repoModel.StatusPaid:
		return model.StatusPaid
	case repoModel.StatusCanceled:
		return model.StatusCanceled
	default:
		return model.StatusPendingPayment
	}
}

// convertModelStatusToRepoStatus конвертирует Status из service в repository model
func convertModelStatusToRepoStatus(modelStatus model.Status) repoModel.Status {
	switch modelStatus {
	case model.StatusPendingPayment:
		return repoModel.StatusPendingPayment
	case model.StatusPaid:
		return repoModel.StatusPaid
	case model.StatusCanceled:
		return repoModel.StatusCanceled
	default:
		return repoModel.StatusPendingPayment
	}
}

// convertModelStatusToOrderStatus конвертирует Status из service model в order API status
func convertModelStatusToOrderStatus(modelStatus model.Status) order_v1.OrderStatus {
	switch modelStatus {
	case model.StatusPendingPayment:
		return order_v1.OrderStatusPENDINGPAYMENT
	case model.StatusPaid:
		return order_v1.OrderStatusPAID
	case model.StatusCanceled:
		return order_v1.OrderStatusCANCELLED
	default:
		return order_v1.OrderStatusPENDINGPAYMENT
	}
}

// convertModelPaymentMethodToOrderPaymentMethod конвертирует PaymentMethod из service model в order API payment method
func convertModelPaymentMethodToOrderPaymentMethod(modelMethod model.PaymentMethod) order_v1.PaymentMethod {
	switch modelMethod {
	case model.PaymentMethodCard:
		return order_v1.PaymentMethodCARD
	case model.PaymentMethodSBP:
		return order_v1.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return order_v1.PaymentMethodCREDITCARD
	case model.PaymentMethodInvestorMoney:
		return order_v1.PaymentMethodINVESTORMONEY
	default:
		return order_v1.PaymentMethodUNKNOWN
	}
}

// convertStringSliceToUUIDSlice конвертирует []string в []uuid.UUID
func convertStringSliceToUUIDSlice(stringSlice []string) []uuid.UUID {
	if stringSlice == nil {
		return nil
	}

	uuidSlice := make([]uuid.UUID, len(stringSlice))
	for i, str := range stringSlice {
		uuidSlice[i] = uuid.MustParse(str)
	}
	return uuidSlice
}

// convertUUIDSliceToStringSlice конвертирует []uuid.UUID в []string
func convertUUIDSliceToStringSlice(uuidSlice []uuid.UUID) []string {
	if uuidSlice == nil {
		return nil
	}

	stringSlice := make([]string, len(uuidSlice))
	for i, u := range uuidSlice {
		stringSlice[i] = u.String()
	}
	return stringSlice
}

// ConvertStringSliceToUUIDSlice конвертирует []string в []uuid.UUID
func ConvertStringSliceToUUIDSlice(stringSlice []string) ([]uuid.UUID, error) {
	if stringSlice == nil {
		return nil, nil
	}

	uuidSlice := make([]uuid.UUID, len(stringSlice))
	for i, str := range stringSlice {
		parsedUUID, err := uuid.Parse(str)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID at index %d: %s", i, str)
		}
		uuidSlice[i] = parsedUUID
	}
	return uuidSlice, nil
}
