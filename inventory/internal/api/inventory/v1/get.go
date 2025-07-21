package v1

import (
	"context"

	"github.com/space-wanderer/microservices/inventory/internal/converter"
	inventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.inventoryService.GetPart(ctx, req.GetUuid())
	if err != nil {
		return nil, err
	}

	grpcPart := converter.ConvertPartToGRPC(part)
	return &inventoryV1.GetPartResponse{
		Part: grpcPart,
	}, nil
}
