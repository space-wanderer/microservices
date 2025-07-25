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

	paymentV1API "github.com/space-wanderer/microservices/payment/internal/api/api/payment/v1"
	"github.com/space-wanderer/microservices/payment/internal/config"
	paymentService "github.com/space-wanderer/microservices/payment/internal/service/payment"
	"github.com/space-wanderer/microservices/shared/pkg/interceptors"
	paymentV1 "github.com/space-wanderer/microservices/shared/pkg/proto/payment/v1"
)

const configPath = "deploy/compose/payment/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %w", err))
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.AppConfig().PaymentGRPC.Address()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Fatalf("failed to close listener: %v", cerr)
		}
	}()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryErrorInterceptor()),
	)

	service := paymentService.NewService()
	api := paymentV1API.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)
	reflection.Register(s)

	go func() {
		log.Printf("gRPS payment listening on %s", config.AppConfig().PaymentGRPC.Address())
		err := s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC payment server...")
	s.GracefulStop()
	log.Println("gRPC payment server stopped")
}
