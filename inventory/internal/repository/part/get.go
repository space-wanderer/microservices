package part

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/space-wanderer/microservices/inventory/internal/model"
	repoModel "github.com/space-wanderer/microservices/inventory/internal/repository/model"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (*repoModel.Part, error) {
	var part repoModel.Part

	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}
		return nil, fmt.Errorf("failed to get part: %w", err)
	}

	return &part, nil
}
