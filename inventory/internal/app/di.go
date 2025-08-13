package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	inventoryV1API "github.com/space-wanderer/microservices/inventory/internal/api/inventory/v1"
	"github.com/space-wanderer/microservices/inventory/internal/config"
	"github.com/space-wanderer/microservices/inventory/internal/repository"
	partRepository "github.com/space-wanderer/microservices/inventory/internal/repository/part"
	"github.com/space-wanderer/microservices/inventory/internal/service"
	partService "github.com/space-wanderer/microservices/inventory/internal/service/part"
	"github.com/space-wanderer/microservices/platform/pkg/closer"
	inventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	inventoryService service.InventoryService

	inventoryRepository repository.InventoryRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.InventoryService(ctx))
	}
	return d.inventoryV1API
}

func (d *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = partService.NewService(d.InventoryRepository(ctx))
	}
	return d.inventoryService
}

func (d *diContainer) InventoryRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = partRepository.NewRepository(ctx, d.MongoDBHandle(ctx))
	}
	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %s\n", err.Error()))
		}
		closer.AddNamed("MongoDB Client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})
		d.mongoDBClient = client
	}
	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.Database())
	}
	return d.mongoDBHandle
}
