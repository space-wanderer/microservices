package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/space-wanderer/microservices/inventory/internal/converter"
	inventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := converter.ConvertFilterFromGRPC(req.GetFilter())
	parts, err := a.inventoryService.ListParts(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal: %s", err.Error())
	}

	grpcParts := make([]*inventoryV1.Part, len(parts))
	for i, part := range parts {
		grpcParts[i] = converter.ConvertPartToGRPC(part)
	}

	return &inventoryV1.ListPartsResponse{Parts: grpcParts}, nil
}
