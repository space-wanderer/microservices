package v1

import (
	generatedInventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
)

type client struct {
	generatedClient generatedInventoryV1.InventoryServiceClient
}

func NewClient(generatedClient generatedInventoryV1.InventoryServiceClient) *client {
	return &client{generatedClient: generatedClient}
}
