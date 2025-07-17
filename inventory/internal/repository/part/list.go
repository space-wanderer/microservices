package part

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	repoModel "github.com/space-wanderer/microservices/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filter *repoModel.PartsFilter) ([]*repoModel.Part, error) {
	// Строим фильтр для MongoDB
	mongoFilter := buildMongoFilter(filter)

	// Выполняем запрос
	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Читаем результаты
	var parts []*repoModel.Part
	if err = cursor.All(ctx, &parts); err != nil {
		return nil, err
	}

	return parts, nil
}

// buildMongoFilter строит фильтр для MongoDB на основе PartsFilter
func buildMongoFilter(filter *repoModel.PartsFilter) bson.M {
	if filter == nil || isEmptyFilter(filter) {
		return bson.M{}
	}

	mongoFilter := bson.M{}

	// Фильтр по UUID
	if len(filter.Uuids) > 0 {
		mongoFilter["uuid"] = bson.M{"$in": filter.Uuids}
	}

	// Фильтр по именам
	if len(filter.Names) > 0 {
		mongoFilter["name"] = bson.M{"$in": filter.Names}
	}

	// Фильтр по категориям
	if len(filter.Categories) > 0 {
		mongoFilter["category"] = bson.M{"$in": filter.Categories}
	}

	// Фильтр по странам производителя
	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	// Фильтр по тегам
	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	return mongoFilter
}

// isEmptyFilter проверяет, пуст ли фильтр
func isEmptyFilter(filter *repoModel.PartsFilter) bool {
	return len(filter.Uuids) == 0 &&
		len(filter.Names) == 0 &&
		len(filter.Categories) == 0 &&
		len(filter.ManufacturerCountries) == 0 &&
		len(filter.Tags) == 0
}
