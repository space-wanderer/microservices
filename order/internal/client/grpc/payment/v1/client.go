package v1

import generatedPaymentV1 "github.com/space-wanderer/microservices/shared/pkg/proto/payment/v1"

type client struct {
	generatedClient generatedPaymentV1.PaymentServiceClient
}

func NewClient(generatedClient generatedPaymentV1.PaymentServiceClient) *client {
	return &client{generatedClient: generatedClient}
}
