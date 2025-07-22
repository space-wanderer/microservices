package service

import (
	"context"

	"github.com/space-wanderer/microservices/inventory/internal/model"
)

type InventoryService interface {
	GetPart(ctx context.Context, uuid string) (*model.Part, error)
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}
