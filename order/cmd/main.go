package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

// Order представляет заказ в памяти
type Order struct {
	OrderUUID       uuid.UUID               `json:"order_uuid"`
	UserUUID        uuid.UUID               `json:"user_uuid"`
	PartUuids       []uuid.UUID             `json:"part_uuids"`
	TotalPrice      float32                 `json:"total_price"`
	TransactionUUID *uuid.UUID              `json:"transaction_uuid,omitempty"`
	PaymentMethod   *order_v1.PaymentMethod `json:"payment_method,omitempty"`
	Status          order_v1.OrderStatus    `json:"status"`
}

// OrderService реализует бизнес-логику заказов
type OrderService struct {
	orders map[uuid.UUID]*Order
	mu     sync.RWMutex

	inventoryClient inventory_v1.InventoryServiceClient
	paymentClient   payment_v1.PaymentServiceClient

	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn
}

// NewOrderService создает новый сервис заказов с gRPC клиентами
func NewOrderService() (*OrderService, error) {
	inventoryConn, err := grpc.NewClient(inventoryServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	paymentConn, err := grpc.NewClient(paymentServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			log.Printf("Ошибка закрытия соединения с inventory: %v", closeErr)
		}
		return nil, err
	}

	return &OrderService{
		orders:          make(map[uuid.UUID]*Order),
		inventoryClient: inventory_v1.NewInventoryServiceClient(inventoryConn),
		paymentClient:   payment_v1.NewPaymentServiceClient(paymentConn),
		inventoryConn:   inventoryConn,
		paymentConn:     paymentConn,
	}, nil
}

// CreateOrder создает новый заказ
func (s *OrderService) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	orderUUID := uuid.New()

	totalPrice, err := s.calculateOrderPrice(ctx, req.PartUuids)
	if err != nil {
		return &order_v1.BadGatewayError{
			Error:   "INVENTORY_SERVICE_ERROR",
			Message: "Ошибка при получении информации о деталях: " + err.Error(),
		}, nil
	}

	order := &Order{
		OrderUUID:  orderUUID,
		UserUUID:   req.UserUUID,
		PartUuids:  req.PartUuids,
		TotalPrice: totalPrice,
		Status:     order_v1.OrderStatusPENDINGPAYMENT,
	}

	s.mu.Lock()
	s.orders[orderUUID] = order
	s.mu.Unlock()

	return &order_v1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

// calculateOrderPrice рассчитывает общую стоимость заказа через InventoryService
func (s *OrderService) calculateOrderPrice(ctx context.Context, partUuids []uuid.UUID) (float32, error) {
	var totalPrice float32

	for _, partUUID := range partUuids {
		req := &inventory_v1.GetPartRequest{
			Uuid: partUUID.String(),
		}

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		resp, err := s.inventoryClient.GetPart(ctx, req)
		if err != nil {
			return 0, err
		}

		totalPrice += float32(resp.Part.Price)
	}

	return totalPrice, nil
}

// GetOrderByUuid получает заказ по UUID
func (s *OrderService) GetOrderByUuid(ctx context.Context, params order_v1.GetOrderByUuidParams) (order_v1.GetOrderByUuidRes, error) {
	s.mu.RLock()
	order, exists := s.orders[params.OrderUUID]
	s.mu.RUnlock()

	if !exists {
		return &order_v1.NotFoundError{
			Code:    404,
			Message: "Заказ не найден",
		}, nil
	}

	// Конвертируем в DTO
	orderDto := order_v1.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
	}

	if order.TransactionUUID != nil {
		orderDto.TransactionUUID = order_v1.NewOptUUID(*order.TransactionUUID)
	}

	if order.PaymentMethod != nil {
		orderDto.PaymentMethod = *order.PaymentMethod
	}

	return &order_v1.GetOrderResponse{
		Order:   orderDto,
		Message: order_v1.NewOptString("Заказ успешно получен"),
	}, nil
}

// PayOrder оплачивает заказ
func (s *OrderService) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {
	s.mu.Lock()
	order, exists := s.orders[params.OrderUUID]
	s.mu.Unlock()

	if !exists {
		return &order_v1.NotFoundError{
			Code:    404,
			Message: "Заказ не найден",
		}, nil
	}

	// Проверяем статус заказа
	if order.Status != order_v1.OrderStatusPENDINGPAYMENT {
		return &order_v1.ForbiddenError{
			Code:    403,
			Message: "Заказ уже оплачен или отменен",
		}, nil
	}

	// Интеграция с PaymentService для обработки платежа
	transactionUUID, err := s.processPayment(ctx, order, req.PaymentMethod)
	if err != nil {
		return &order_v1.BadGatewayError{
			Error:   "PAYMENT_SERVICE_ERROR",
			Message: "Ошибка при обработке платежа: " + err.Error(),
		}, nil
	}

	// Обновляем заказ
	s.mu.Lock()
	order.Status = order_v1.OrderStatusPAID
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &req.PaymentMethod
	s.mu.Unlock()

	return &order_v1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

// processPayment обрабатывает платеж через PaymentService
func (s *OrderService) processPayment(ctx context.Context, order *Order, paymentMethod order_v1.PaymentMethod) (uuid.UUID, error) {
	var grpcPaymentMethod payment_v1.PaymentMethod
	switch paymentMethod {
	case order_v1.PaymentMethodCARD:
		grpcPaymentMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case order_v1.PaymentMethodSBP:
		grpcPaymentMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_SPB
	case order_v1.PaymentMethodCREDITCARD:
		grpcPaymentMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case order_v1.PaymentMethodINVESTORMONEY:
		grpcPaymentMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		grpcPaymentMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	req := &payment_v1.PayOrderRequest{
		OrderUuid:     order.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: grpcPaymentMethod,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := s.paymentClient.PayOrder(ctx, req)
	if err != nil {
		return uuid.Nil, err
	}

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return uuid.Nil, err
	}

	return transactionUUID, nil
}

// CancelOrderByUuid отменяет заказ
func (s *OrderService) CancelOrderByUuid(ctx context.Context, params order_v1.CancelOrderByUuidParams) (order_v1.CancelOrderByUuidRes, error) {
	s.mu.Lock()
	order, exists := s.orders[params.OrderUUID]
	s.mu.Unlock()

	if !exists {
		return &order_v1.NotFoundError{
			Code:    404,
			Message: "Заказ не найден",
		}, nil
	}

	// Проверяем статус заказа
	if order.Status == order_v1.OrderStatusPAID {
		return &order_v1.ConflictError{
			Code:    409,
			Message: "Заказ уже оплачен и не может быть отменен",
		}, nil
	}

	// Отменяем заказ
	s.mu.Lock()
	order.Status = order_v1.OrderStatusCANCELLED
	s.mu.Unlock()

	return &order_v1.CancelOrderByUuidNoContent{}, nil
}

// NewError создает новую ошибку
func (s *OrderService) NewError(ctx context.Context, err error) *order_v1.GenericErrorStatusCode {
	return &order_v1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: order_v1.GenericError{
			Message: err.Error(),
		},
	}
}

// Close закрывает gRPC соединения
func (s *OrderService) Close() {
	log.Println("Закрытие gRPC соединений...")
	if err := s.inventoryConn.Close(); err != nil {
		log.Printf("Ошибка закрытия соединения с inventory: %v", err)
	}
	if err := s.paymentConn.Close(); err != nil {
		log.Printf("Ошибка закрытия соединения с payment: %v", err)
	}
}

func main() {
	// Создаем сервис заказов
	orderService, err := NewOrderService()
	if err != nil {
		log.Printf("❌ Ошибка создания сервиса заказов: %v", err)
		os.Exit(1)
	}

	// Создаем OpenAPI сервер
	s, err := order_v1.NewServer(orderService)
	if err != nil {
		log.Printf("❌ Ошибка создания сервера: %v", err)
		orderService.Close()
		os.Exit(1)
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

	orderService.Close()

	log.Println("✅ Сервер остановлен")
}
