package converter

import (
	"time"

	"github.com/space-wanderer/microservices/inventory/internal/repository/model"
	inventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertProtoPartToRepoPart конвертирует Part из proto в repository model
func ConvertProtoPartToRepoPart(protoPart *inventoryV1.Part) *model.Part {
	if protoPart == nil {
		return nil
	}

	var dimensions *model.Dimensions
	if protoPart.Dimensions != nil {
		dimensions = &model.Dimensions{
			Length: protoPart.Dimensions.Length,
			Width:  protoPart.Dimensions.Width,
			Height: protoPart.Dimensions.Height,
			Weight: protoPart.Dimensions.Weight,
		}
	}

	var manufacturer *model.Manufacturer
	if protoPart.Manufacturer != nil {
		manufacturer = &model.Manufacturer{
			Name:    protoPart.Manufacturer.Name,
			Country: protoPart.Manufacturer.Country,
			Website: protoPart.Manufacturer.Website,
		}
	}

	return &model.Part{
		UUID:          protoPart.Uuid,
		Name:          protoPart.Name,
		Description:   protoPart.Description,
		Price:         protoPart.Price,
		StockQuantity: protoPart.StockQuantity,
		Category:      convertProtoCategoryToRepoCategory(protoPart.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          protoPart.Tags,
		Metadata:      convertProtoMetadataToRepoMetadata(protoPart.Metadata),
		CreatedAt:     convertProtoTimestampToTime(protoPart.CreatedAt),
		UpdatedAt:     convertProtoTimestampToTime(protoPart.UpdatedAt),
	}
}

// ConvertRepoPartToProtoPart конвертирует Part из repository model в proto
func ConvertRepoPartToProtoPart(repoPart *model.Part) *inventoryV1.Part {
	if repoPart == nil {
		return nil
	}

	var dimensions *inventoryV1.Dimensions
	if repoPart.Dimensions != nil {
		dimensions = &inventoryV1.Dimensions{
			Length: repoPart.Dimensions.Length,
			Width:  repoPart.Dimensions.Width,
			Height: repoPart.Dimensions.Height,
			Weight: repoPart.Dimensions.Weight,
		}
	}

	var manufacturer *inventoryV1.Manufacturer
	if repoPart.Manufacturer != nil {
		manufacturer = &inventoryV1.Manufacturer{
			Name:    repoPart.Manufacturer.Name,
			Country: repoPart.Manufacturer.Country,
			Website: repoPart.Manufacturer.Website,
		}
	}

	return &inventoryV1.Part{
		Uuid:          repoPart.UUID,
		Name:          repoPart.Name,
		Description:   repoPart.Description,
		Price:         repoPart.Price,
		StockQuantity: repoPart.StockQuantity,
		Category:      convertRepoCategoryToProtoCategory(repoPart.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          repoPart.Tags,
		Metadata:      convertRepoMetadataToProtoMetadata(repoPart.Metadata),
		CreatedAt:     timestamppb.New(repoPart.CreatedAt),
		UpdatedAt:     timestamppb.New(repoPart.UpdatedAt),
	}
}

func convertProtoCategoryToRepoCategory(protoCategory inventoryV1.Category) model.Category {
	switch protoCategory {
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

func convertRepoCategoryToProtoCategory(repoCategory model.Category) inventoryV1.Category {
	switch repoCategory {
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

func convertProtoMetadataToRepoMetadata(protoMetadata map[string]*inventoryV1.Value) map[string]*model.Value {
	if protoMetadata == nil {
		return nil
	}

	repoMetadata := make(map[string]*model.Value)
	for key, protoValue := range protoMetadata {
		if protoValue != nil {
			repoMetadata[key] = convertProtoValueToRepoValue(protoValue)
		}
	}
	return repoMetadata
}

func convertRepoMetadataToProtoMetadata(repoMetadata map[string]*model.Value) map[string]*inventoryV1.Value {
	if repoMetadata == nil {
		return nil
	}

	protoMetadata := make(map[string]*inventoryV1.Value)
	for key, repoValue := range repoMetadata {
		if repoValue != nil {
			protoMetadata[key] = convertRepoValueToProtoValue(repoValue)
		}
	}
	return protoMetadata
}

func convertProtoValueToRepoValue(protoValue *inventoryV1.Value) *model.Value {
	if protoValue == nil {
		return nil
	}

	var value model.Value
	switch v := protoValue.Value.(type) {
	case *inventoryV1.Value_StringValue:
		value = &model.StringValue{StringValue: v.StringValue}
	case *inventoryV1.Value_Int64Value:
		value = &model.Int64Value{Int64Value: v.Int64Value}
	case *inventoryV1.Value_DoubleValue:
		value = &model.DoubleValue{DoubleValue: v.DoubleValue}
	case *inventoryV1.Value_BoolValue:
		value = &model.BoolValue{BoolValue: v.BoolValue}
	default:
		return nil
	}
	return &value
}

func convertRepoValueToProtoValue(repoValue *model.Value) *inventoryV1.Value {
	if repoValue == nil {
		return nil
	}

	switch v := (*repoValue).(type) {
	case *model.StringValue:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_StringValue{StringValue: v.StringValue},
		}
	case *model.Int64Value:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_Int64Value{Int64Value: v.Int64Value},
		}
	case *model.DoubleValue:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_DoubleValue{DoubleValue: v.DoubleValue},
		}
	case *model.BoolValue:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_BoolValue{BoolValue: v.BoolValue},
		}
	default:
		return nil
	}
}

func convertProtoTimestampToTime(protoTimestamp *timestamppb.Timestamp) time.Time {
	if protoTimestamp == nil {
		return time.Time{}
	}
	return protoTimestamp.AsTime()
}
