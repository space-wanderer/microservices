package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1API "github.com/space-wanderer/microservices/inventory/internal/api/inventory/v1"
	inventoryRepository "github.com/space-wanderer/microservices/inventory/internal/repository/part"
	inventoryService "github.com/space-wanderer/microservices/inventory/internal/service/part"
	sharedInterceptors "github.com/space-wanderer/microservices/shared/pkg/interceptors"
	inventoryV1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
)

const grpsPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpsPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Fatalf("failed to close listener: %v", cerr)
		}
	}()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(sharedInterceptors.UnaryErrorInterceptor()),
	)

	repo := inventoryRepository.NewRepository()
	service := inventoryService.NewService(repo)
	api := inventoryV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)
	reflection.Register(s)

	go func() {
		log.Printf("gRPS inventory listening on %d", grpsPort)
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
