package v1

import (
	"context"

	"github.com/space-wanderer/microservices/inventory/internal/converter"
	inventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.inventoryService.GetPart(ctx, req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "NotFound: %s", req.GetUuid())
	}

	grpcPart := converter.ConvertPartToGRPC(part)
	return &inventoryV1.GetPartResponse{
		Part: grpcPart,
	}, nil
}
