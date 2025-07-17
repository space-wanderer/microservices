package converter

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	serviceModel "github.com/space-wanderer/microservices/inventory/internal/model"
	repoModel "github.com/space-wanderer/microservices/inventory/internal/repository/model"
)

// ConvertServicePartToRepoPart конвертирует Part из service model в repository model
func ConvertServicePartToRepoPart(servicePart *serviceModel.Part) *repoModel.Part {
	if servicePart == nil {
		return nil
	}

	var dimensions *repoModel.Dimensions
	if servicePart.Dimensions != nil {
		dimensions = &repoModel.Dimensions{
			Length: servicePart.Dimensions.Length,
			Width:  servicePart.Dimensions.Width,
			Height: servicePart.Dimensions.Height,
			Weight: servicePart.Dimensions.Weight,
		}
	}

	var manufacturer *repoModel.Manufacturer
	if servicePart.Manufacturer != nil {
		manufacturer = &repoModel.Manufacturer{
			Name:    servicePart.Manufacturer.Name,
			Country: servicePart.Manufacturer.Country,
			Website: servicePart.Manufacturer.Website,
		}
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	return &repoModel.Part{
		UUID:          servicePart.UUID,
		Name:          servicePart.Name,
		Description:   servicePart.Description,
		Price:         servicePart.Price,
		StockQuantity: servicePart.StockQuantity,
		Category:      convertServiceCategoryToRepoCategory(servicePart.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          servicePart.Tags,
		Metadata:      convertServiceMetadataToRepoMetadata(servicePart.Metadata),
		CreatedAt:     convertServiceTimeToPrimitiveDateTime(servicePart.CreatedAt, now),
		UpdatedAt:     now, // Всегда обновляем время изменения
	}
}

// ConvertRepoPartToServicePart конвертирует Part из repository model в service model
func ConvertRepoPartToServicePart(repoPart *repoModel.Part) *serviceModel.Part {
	if repoPart == nil {
		return nil
	}

	var dimensions *serviceModel.Dimensions
	if repoPart.Dimensions != nil {
		dimensions = &serviceModel.Dimensions{
			Length: repoPart.Dimensions.Length,
			Width:  repoPart.Dimensions.Width,
			Height: repoPart.Dimensions.Height,
			Weight: repoPart.Dimensions.Weight,
		}
	}

	var manufacturer *serviceModel.Manufacturer
	if repoPart.Manufacturer != nil {
		manufacturer = &serviceModel.Manufacturer{
			Name:    repoPart.Manufacturer.Name,
			Country: repoPart.Manufacturer.Country,
			Website: repoPart.Manufacturer.Website,
		}
	}

	return &serviceModel.Part{
		UUID:          repoPart.UUID,
		Name:          repoPart.Name,
		Description:   repoPart.Description,
		Price:         repoPart.Price,
		StockQuantity: repoPart.StockQuantity,
		Category:      convertRepoCategoryToServiceCategory(repoPart.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          repoPart.Tags,
		Metadata:      convertRepoMetadataToServiceMetadata(repoPart.Metadata),
		CreatedAt:     convertPrimitiveDateTimeToServiceTime(repoPart.CreatedAt),
		UpdatedAt:     convertPrimitiveDateTimeToServiceTime(repoPart.UpdatedAt),
	}
}

// ConvertServiceFilterToRepoFilter конвертирует PartsFilter из service model в repository model
func ConvertServiceFilterToRepoFilter(serviceFilter *serviceModel.PartsFilter) *repoModel.PartsFilter {
	if serviceFilter == nil {
		return nil
	}

	var categories []repoModel.Category
	if serviceFilter.Categories != nil {
		categories = make([]repoModel.Category, len(serviceFilter.Categories))
		for i, category := range serviceFilter.Categories {
			categories[i] = convertServiceCategoryToRepoCategory(category)
		}
	}

	return &repoModel.PartsFilter{
		Uuids:                 serviceFilter.Uuids,
		Names:                 serviceFilter.Names,
		Categories:            categories,
		ManufacturerCountries: serviceFilter.ManufacturerCountries,
		Tags:                  serviceFilter.Tags,
	}
}

// Конвертеры для категорий
func convertServiceCategoryToRepoCategory(serviceCategory serviceModel.Category) repoModel.Category {
	switch serviceCategory {
	case serviceModel.CategoryEngine:
		return repoModel.CategoryEngine
	case serviceModel.CategoryFuel:
		return repoModel.CategoryFuel
	case serviceModel.CategoryPorthole:
		return repoModel.CategoryPorthole
	case serviceModel.CategoryWing:
		return repoModel.CategoryWing
	default:
		return repoModel.CategoryUnknown
	}
}

func convertRepoCategoryToServiceCategory(repoCategory repoModel.Category) serviceModel.Category {
	switch repoCategory {
	case repoModel.CategoryEngine:
		return serviceModel.CategoryEngine
	case repoModel.CategoryFuel:
		return serviceModel.CategoryFuel
	case repoModel.CategoryPorthole:
		return serviceModel.CategoryPorthole
	case repoModel.CategoryWing:
		return serviceModel.CategoryWing
	default:
		return serviceModel.CategoryUnknown
	}
}

// Конвертеры для метаданных
func convertServiceMetadataToRepoMetadata(serviceMetadata map[string]*serviceModel.Value) map[string]*repoModel.Value {
	if serviceMetadata == nil {
		return nil
	}

	repoMetadata := make(map[string]*repoModel.Value)
	for key, serviceValue := range serviceMetadata {
		if serviceValue != nil {
			repoMetadata[key] = convertServiceValueToRepoValue(serviceValue)
		}
	}
	return repoMetadata
}

func convertRepoMetadataToServiceMetadata(repoMetadata map[string]*repoModel.Value) map[string]*serviceModel.Value {
	if repoMetadata == nil {
		return nil
	}

	serviceMetadata := make(map[string]*serviceModel.Value)
	for key, repoValue := range repoMetadata {
		if repoValue != nil {
			serviceMetadata[key] = convertRepoValueToServiceValue(repoValue)
		}
	}
	return serviceMetadata
}

// Конвертеры для Value типов
func convertServiceValueToRepoValue(serviceValue *serviceModel.Value) *repoModel.Value {
	if serviceValue == nil {
		return nil
	}

	var value repoModel.Value
	switch v := (*serviceValue).(type) {
	case *serviceModel.StringValue:
		value = &repoModel.StringValue{StringValue: v.StringValue}
	case *serviceModel.Int64Value:
		value = &repoModel.Int64Value{Int64Value: v.Int64Value}
	case *serviceModel.DoubleValue:
		value = &repoModel.DoubleValue{DoubleValue: v.DoubleValue}
	case *serviceModel.BoolValue:
		value = &repoModel.BoolValue{BoolValue: v.BoolValue}
	default:
		return nil
	}
	return &value
}

func convertRepoValueToServiceValue(repoValue *repoModel.Value) *serviceModel.Value {
	if repoValue == nil {
		return nil
	}

	var value serviceModel.Value
	switch v := (*repoValue).(type) {
	case *repoModel.StringValue:
		value = &serviceModel.StringValue{StringValue: v.StringValue}
	case *repoModel.Int64Value:
		value = &serviceModel.Int64Value{Int64Value: v.Int64Value}
	case *repoModel.DoubleValue:
		value = &serviceModel.DoubleValue{DoubleValue: v.DoubleValue}
	case *repoModel.BoolValue:
		value = &serviceModel.BoolValue{BoolValue: v.BoolValue}
	default:
		return nil
	}
	return &value
}

// Конвертеры для времени
func convertServiceTimeToPrimitiveDateTime(serviceTime time.Time, fallback primitive.DateTime) primitive.DateTime {
	if serviceTime.IsZero() {
		return fallback
	}
	return primitive.NewDateTimeFromTime(serviceTime)
}

func convertPrimitiveDateTimeToServiceTime(repoTime primitive.DateTime) time.Time {
	return repoTime.Time()
}
