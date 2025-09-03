package v1

import (
	"context"

	"github.com/space-wanderer/microservices/order/internal/client/converter"
	"github.com/space-wanderer/microservices/order/internal/model"
	grpcAuth "github.com/space-wanderer/microservices/platform/pkg/middleware/grpc"
	genaratedInventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]*model.Part, error) {
	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)
	parts, err := c.generatedClient.ListParts(ctx, &genaratedInventoryV1.ListPartsRequest{
		Filter: converter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}
	return converter.PartListProtoToModel(parts.Parts), nil
}
