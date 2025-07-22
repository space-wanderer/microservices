package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/space-wanderer/microservices/order/internal/api/order/v1"
	grpcClient "github.com/space-wanderer/microservices/order/internal/client/grpc"
	orderRepository "github.com/space-wanderer/microservices/order/internal/repository/order"
	orderService "github.com/space-wanderer/microservices/order/internal/service/order"
	order_v1 "github.com/space-wanderer/microservices/shared/pkg/api/order/v1"
	inventory_v1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/space-wanderer/microservices/shared/pkg/proto/payment/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second

	// Адреса gRPC сервисов
	inventoryServiceAddr = "localhost:50051"
	paymentServiceAddr   = "localhost:50052"
)

func main() {
	// Создаем gRPC соединения
	inventoryConn, err := grpc.NewClient(inventoryServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("❌ Ошибка подключения к inventory service: %v", err)
		return
	}
	defer func() {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			log.Printf("❌ Ошибка закрытия соединения с inventory: %v", closeErr)
		}
	}()

	paymentConn, err := grpc.NewClient(paymentServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("❌ Ошибка подключения к payment service: %v", err)
		return
	}
	defer func() {
		if closeErr := paymentConn.Close(); closeErr != nil {
			log.Printf("❌ Ошибка закрытия соединения с payment: %v", closeErr)
		}
	}()

	// Создаем адаптированные клиенты
	inventoryClient := inventory_v1.NewInventoryServiceClient(inventoryConn)
	paymentClient := payment_v1.NewPaymentServiceClient(paymentConn)

	// Создаем адаптеры
	inventoryAdapter := grpcClient.NewInventoryClient(inventoryClient)
	paymentAdapter := grpcClient.NewPaymentClient(paymentClient)

	// Создаем репозиторий и сервис
	repo := orderRepository.NewRepository()
	service := orderService.NewOrderService(repo, inventoryAdapter, paymentAdapter)
	api := orderV1API.NewAPI(service)

	// Создаем OpenAPI сервер
	s, err := order_v1.NewServer(api)
	if err != nil {
		log.Printf("❌ Ошибка создания сервера: %v", err)
		return
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", s)

	// Запускаем HTTP-сервер
	httpServer := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		log.Printf("🔗 Интеграция с InventoryService: %s\n", inventoryServiceAddr)
		log.Printf("💳 Интеграция с PaymentService: %s\n", paymentServiceAddr)
		err = httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
