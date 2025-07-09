package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentV1 "github.com/space-wanderer/microservices/shared/pkg/proto/payment/v1"
)

const grpsPort = 50052

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func (s *paymentService) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transaction := &paymentV1.PayOrderResponse{
		TransactionUuid: uuid.New().String(),
	}

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transaction.TransactionUuid)

	return transaction, nil
}

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

	service := &paymentService{}

	paymentV1.RegisterPaymentServiceServer(s, service)
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
