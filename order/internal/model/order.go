package model

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUuids       []string
	TotalPrice      float32
	TransactionUUID *string
	PaymentMethod   PaymentMethod
	Status          Status
}

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type Status string

const (
	StatusPendingPayment Status = "PENDING_PAYMENT"
	StatusPaid           Status = "PAID"
	StatusCanceled       Status = "CANCELED"
	StatusAssembled      Status = "ASSEMBLED"
)
