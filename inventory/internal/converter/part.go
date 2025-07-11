package converter

import (
	"github.com/space-wanderer/microservices/inventory/internal/model"
	repoModel "github.com/space-wanderer/microservices/inventory/internal/repository/model"
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
