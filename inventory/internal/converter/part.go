package converter

import (
	"github.com/space-wanderer/microservices/inventory/internal/model"
	repoModel "github.com/space-wanderer/microservices/inventory/internal/repository/model"
	inventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertRepoPartToModelPart конвертирует Part из repository model в service model
func ConvertRepoPartToModelPart(repoPart *repoModel.Part) *model.Part {
	if repoPart == nil {
		return nil
	}

	var dimensions *model.Dimensions
	if repoPart.Dimensions != nil {
		dimensions = &model.Dimensions{
			Length: repoPart.Dimensions.Length,
			Width:  repoPart.Dimensions.Width,
			Height: repoPart.Dimensions.Height,
			Weight: repoPart.Dimensions.Weight,
		}
	}

	var manufacturer *model.Manufacturer
	if repoPart.Manufacturer != nil {
		manufacturer = &model.Manufacturer{
			Name:    repoPart.Manufacturer.Name,
			Country: repoPart.Manufacturer.Country,
			Website: repoPart.Manufacturer.Website,
		}
	}

	return &model.Part{
		UUID:          repoPart.UUID,
		Name:          repoPart.Name,
		Description:   repoPart.Description,
		Price:         repoPart.Price,
		StockQuantity: repoPart.StockQuantity,
		Category:      model.Category(repoPart.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          repoPart.Tags,
		Metadata:      convertRepoMetadataToModelMetadata(repoPart.Metadata),
		CreatedAt:     repoPart.CreatedAt,
		UpdatedAt:     repoPart.UpdatedAt,
	}
}

func convertRepoMetadataToModelMetadata(repoMetadata map[string]*repoModel.Value) map[string]*model.Value {
	if repoMetadata == nil {
		return nil
	}

	modelMetadata := make(map[string]*model.Value)
	for key, repoValue := range repoMetadata {
		if repoValue != nil {
			value := convertRepoValueToModelValue(repoValue)
			modelMetadata[key] = &value
		}
	}
	return modelMetadata
}

func convertRepoValueToModelValue(repoValue *repoModel.Value) model.Value {
	if repoValue == nil {
		return nil
	}

	switch v := (*repoValue).(type) {
	case *repoModel.StringValue:
		return &model.StringValue{StringValue: v.StringValue}
	case *repoModel.Int64Value:
		return &model.Int64Value{Int64Value: v.Int64Value}
	case *repoModel.DoubleValue:
		return &model.DoubleValue{DoubleValue: v.DoubleValue}
	case *repoModel.BoolValue:
		return &model.BoolValue{BoolValue: v.BoolValue}
	default:
		return nil
	}
}

// ConvertModelPartsFilterToRepoPartsFilter конвертирует PartsFilter из service model в repository model
func ConvertModelPartsFilterToRepoPartsFilter(modelFilter *model.PartsFilter) *repoModel.PartsFilter {
	if modelFilter == nil {
		return nil
	}

	return &repoModel.PartsFilter{
		Uuids:                 modelFilter.Uuids,
		Names:                 modelFilter.Names,
		Categories:            convertModelCategoriesToRepoCategories(modelFilter.Categories),
		ManufacturerCountries: modelFilter.ManufacturerCountries,
		Tags:                  modelFilter.Tags,
	}
}

func convertModelCategoriesToRepoCategories(modelCategories []model.Category) []repoModel.Category {
	if modelCategories == nil {
		return nil
	}

	repoCategories := make([]repoModel.Category, len(modelCategories))
	for i, category := range modelCategories {
		repoCategories[i] = repoModel.Category(category)
	}
	return repoCategories
}

// ConvertRepoPartsToModelParts конвертирует массив Parts из repository model в service model
func ConvertRepoPartsToModelParts(repoParts []*repoModel.Part) []*model.Part {
	if repoParts == nil {
		return nil
	}

	modelParts := make([]*model.Part, len(repoParts))
	for i, repoPart := range repoParts {
		modelParts[i] = ConvertRepoPartToModelPart(repoPart)
	}
	return modelParts
}

// convertPartToGRPC конвертирует внутреннюю модель Part в gRPC модель
func ConvertPartToGRPC(part *model.Part) *inventoryV1.Part {
	grpcPart := &inventoryV1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      convertCategoryToGRPC(part.Category),
		Tags:          part.Tags,
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     timestamppb.New(part.UpdatedAt),
	}

	if part.Dimensions != nil {
		grpcPart.Dimensions = &inventoryV1.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		}
	}

	if part.Manufacturer != nil {
		grpcPart.Manufacturer = &inventoryV1.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		}
	}

	return grpcPart
}

// convertCategoryToGRPC конвертирует внутреннюю категорию в gRPC категорию
func convertCategoryToGRPC(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}

// ConvertFilterFromGRPC конвертирует gRPC фильтр в модель фильтра
func ConvertFilterFromGRPC(grpcFilter *inventoryV1.PartsFilter) *model.PartsFilter {
	if grpcFilter == nil {
		return &model.PartsFilter{}
	}
	return &model.PartsFilter{
		Uuids:                 grpcFilter.GetUuids(),
		Names:                 grpcFilter.GetNames(),
		Categories:            convertGRPCCategoriesToModelCategories(grpcFilter.GetCategories()),
		ManufacturerCountries: grpcFilter.GetManufacturerCountries(),
		Tags:                  grpcFilter.GetTags(),
	}
}

// convertGRPCCategoriesToModelCategories конвертирует gRPC категории в модель категории
func convertGRPCCategoriesToModelCategories(grpcCategories []inventoryV1.Category) []model.Category {
	if grpcCategories == nil {
		return nil
	}

	modelCategories := make([]model.Category, len(grpcCategories))
	for i, category := range grpcCategories {
		modelCategories[i] = convertGRPCCategoryToModelCategory(category)
	}
	return modelCategories
}

// convertGRPCCategoryToModelCategory конвертирует gRPC категорию в модель категории
func convertGRPCCategoryToModelCategory(grpcCategory inventoryV1.Category) model.Category {
	switch grpcCategory {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}
