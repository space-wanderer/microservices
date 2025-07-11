package part

import (
	"context"
	"fmt"

	repoModel "github.com/space-wanderer/microservices/inventory/internal/repository/model"
)

func (r *repository) GetPart(_ context.Context, uuid string) (*repoModel.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[uuid]
	if !ok {
		return nil, fmt.Errorf("part not found: %s", uuid)
	}
	return part, nil
}
