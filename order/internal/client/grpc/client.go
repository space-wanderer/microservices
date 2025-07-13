package grpc

import (
	"context"

	inventoryV1 "github.com/space-wanderer/microservices/order/internal/client/grpc/inventory/v1"
	paymentV1 "github.com/space-wanderer/microservices/order/internal/client/grpc/payment/v1"
	"github.com/space-wanderer/microservices/order/internal/model"
	genaratedInventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
	generatedPaymentV1 "github.com/space-wanderer/microservices/shared/pkg/proto/payment/v1"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]*model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error)
}

func NewInventoryClient(generatedClient genaratedInventoryV1.InventoryServiceClient) InventoryClient {
	return inventoryV1.NewClient(generatedClient)
}

func NewPaymentClient(generatedClient generatedPaymentV1.PaymentServiceClient) PaymentClient {
	return paymentV1.NewClient(generatedClient)
}
