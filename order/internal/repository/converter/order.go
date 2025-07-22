package converter

import (
	"github.com/space-wanderer/microservices/order/internal/model"
	repoModel "github.com/space-wanderer/microservices/order/internal/repository/model"
)

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
