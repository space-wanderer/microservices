package converter

import (
	"github.com/space-wanderer/microservices/order/internal/model"
	genaratedInventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
)

// PartsFilterToProto конвертирует PartsFilter из модели в proto
func PartsFilterToProto(filter model.PartsFilter) *genaratedInventoryV1.PartsFilter {
	return &genaratedInventoryV1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            convertCategoriesToProto(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

// PartListProtoToModel конвертирует список Parts из proto в модель
func PartListProtoToModel(protoParts []*genaratedInventoryV1.Part) []*model.Part {
	if protoParts == nil {
		return nil
	}

	parts := make([]*model.Part, len(protoParts))
	for i, protoPart := range protoParts {
		parts[i] = PartProtoToModel(protoPart)
	}
	return parts
}

// PartProtoToModel конвертирует Part из proto в модель
func PartProtoToModel(protoPart *genaratedInventoryV1.Part) *model.Part {
	if protoPart == nil {
		return nil
	}

	return &model.Part{
		UUID:  protoPart.Uuid,
		Name:  protoPart.Name,
		Price: protoPart.Price,
	}
}

// convertCategoriesToProto конвертирует категории из модели в proto
func convertCategoriesToProto(categories []model.Category) []genaratedInventoryV1.Category {
	if categories == nil {
		return nil
	}

	protoCategories := make([]genaratedInventoryV1.Category, len(categories))
	for i, category := range categories {
		protoCategories[i] = convertCategoryToProto(category)
	}
	return protoCategories
}

// convertCategoryToProto конвертирует категорию из модели в proto
func convertCategoryToProto(category model.Category) genaratedInventoryV1.Category {
	switch category {
	case model.CategoryEngine:
		return genaratedInventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return genaratedInventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return genaratedInventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return genaratedInventoryV1.Category_CATEGORY_WING
	default:
		return genaratedInventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}
