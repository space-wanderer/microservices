package part

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	repoModel "github.com/space-wanderer/microservices/inventory/internal/repository/model"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (*repoModel.Part, error) {
	var part repoModel.Part

	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("part not found: %s", uuid)
		}
		return nil, fmt.Errorf("failed to get part: %w", err)
	}

	return &part, nil
}
