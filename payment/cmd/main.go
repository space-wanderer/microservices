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
	paymentService "github.com/space-wanderer/microservices/payment/internal/service/payment"
	paymentV1 "github.com/space-wanderer/microservices/shared/pkg/proto/payment/v1"
)

const grpsPort = 50052

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

	s := grpc.NewServer()

	myPaymentService := &paymentService.Service{}
	service := paymentService.NewService(myPaymentService)
	api := paymentV1API.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)
	reflection.Register(s)

	go func() {
		log.Printf("gRPS payment listening on %d", grpsPort)
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
