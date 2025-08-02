package v1

import (
	"context"

	orderV1 "github.com/space-wanderer/microservices/shared/pkg/api/order/v1"
)

func (a *api) CancelOrderByUuid(ctx context.Context, params orderV1.CancelOrderByUuidParams) (orderV1.CancelOrderByUuidRes, error) {
	_, err := a.orderService.CancelOrderByUuid(ctx, params.OrderUUID.String())
	if err != nil {
		return nil, err
	}

	return &orderV1.CancelOrderByUuidNoContent{}, nil
}
