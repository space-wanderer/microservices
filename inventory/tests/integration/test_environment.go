package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
)

// InsertTestPart — вставляет тестовую деталь в коллекцию Mongo и возвращает её UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now()

	partDoc := bson.M{
		"_id":            partUUID,
		"name":           gofakeit.ProductName(),
		"description":    gofakeit.Sentence(15),
		"price":          gofakeit.Float64Range(1000, 500000),
		"stock_quantity": int64(gofakeit.Number(1, 100)),
		"category":       "ENGINE",
		"dimensions": bson.M{
			"length": gofakeit.Float64Range(10, 500),
			"width":  gofakeit.Float64Range(10, 200),
			"height": gofakeit.Float64Range(5, 100),
			"weight": gofakeit.Float64Range(1, 1000),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": "https://" + gofakeit.DomainName(),
		},
		"tags":       []string{gofakeit.Word(), gofakeit.Word(), gofakeit.Word()},
		"created_at": primitive.NewDateTimeFromTime(now),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// InsertTestPartWithData — вставляет тестовую деталь с заданными данными
func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, part *inventoryV1.Part) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now()

	partDoc := bson.M{
		"_id":            partUUID,
		"name":           part.GetName(),
		"description":    part.GetDescription(),
		"price":          part.GetPrice(),
		"stock_quantity": part.GetStockQuantity(),
		"category":       part.GetCategory().String(),
		"dimensions": bson.M{
			"length": part.GetDimensions().GetLength(),
			"width":  part.GetDimensions().GetWidth(),
			"height": part.GetDimensions().GetHeight(),
			"weight": part.GetDimensions().GetWeight(),
		},
		"manufacturer": bson.M{
			"name":    part.GetManufacturer().GetName(),
			"country": part.GetManufacturer().GetCountry(),
			"website": part.GetManufacturer().GetWebsite(),
		},
		"tags":       part.GetTags(),
		"created_at": primitive.NewDateTimeFromTime(now),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// GetTestPartInfo — возвращает тестовую информацию о детали
func (env *TestEnvironment) GetTestPartInfo() *inventoryV1.Part {
	return &inventoryV1.Part{
		Name:          "Ионный двигатель X-2000",
		Description:   "Высокоэффективный ионный двигатель для межпланетных полетов",
		Price:         150000.0,
		StockQuantity: 5,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 120.0,
			Width:  80.0,
			Height: 60.0,
			Weight: 250.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "КосмоТех",
			Country: "Россия",
			Website: "https://cosmotech.ru",
		},
		Tags:      []string{"ионный", "двигатель", "межпланетный", "высокоэффективный"},
		CreatedAt: timestamppb.New(time.Now().Add(-30 * 24 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now()),
	}
}

// GetUpdatedPartInfo — возвращает обновленную информацию о детали
func (env *TestEnvironment) GetUpdatedPartInfo() *inventoryV1.Part {
	return &inventoryV1.Part{
		Name:          "Плазменный двигатель P-500",
		Description:   "Мощный плазменный двигатель для тяжелых грузов",
		Price:         200000.0,
		StockQuantity: 3,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 150.0,
			Width:  100.0,
			Height: 80.0,
			Weight: 400.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "StarTech Industries",
			Country: "США",
			Website: "https://startech.com",
		},
		Tags:      []string{"плазменный", "двигатель", "тяжелый", "грузовой"},
		CreatedAt: timestamppb.New(time.Now().Add(-45 * 24 * time.Hour)),
		UpdatedAt: timestamppb.New(time.Now()),
	}
}

// ClearPartsCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
