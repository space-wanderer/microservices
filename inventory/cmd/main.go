package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1API "github.com/space-wanderer/microservices/inventory/internal/api/inventory/v1"
	"github.com/space-wanderer/microservices/inventory/internal/config"
	inventoryRepository "github.com/space-wanderer/microservices/inventory/internal/repository/part"
	inventoryService "github.com/space-wanderer/microservices/inventory/internal/service/part"
	"github.com/space-wanderer/microservices/shared/pkg/interceptors"
	inventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
)

const configPath = "deploy/compose/inventory/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %w", err))
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.AppConfig().InventoryGRPC.Address()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Fatalf("failed to close listener: %v", cerr)
		}
	}()

	// Подключение к MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		config.AppConfig().Mongo.URI(),
	))
	if err != nil {
		log.Printf("failed to connect to MongoDB: %v", err)
		return
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Проверка подключения
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping MongoDB: %v", err)
		return
	}

	db := client.Database(config.AppConfig().Mongo.Database())

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryErrorInterceptor()),
	)

	repo := inventoryRepository.NewRepository(db)
	service := inventoryService.NewService(repo)
	api := inventoryV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)
	reflection.Register(s)

	go func() {
		log.Printf("gRPS inventory listening on %s", config.AppConfig().InventoryGRPC.Address())
		err := s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("gRPC server stopped")
}
