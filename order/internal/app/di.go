package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/space-wanderer/microservices/order/internal/api/order/v1"
	grpcClient "github.com/space-wanderer/microservices/order/internal/client/grpc"
	"github.com/space-wanderer/microservices/order/internal/config"
	"github.com/space-wanderer/microservices/order/internal/repository"
	orderRepository "github.com/space-wanderer/microservices/order/internal/repository/order"
	"github.com/space-wanderer/microservices/order/internal/service"
	orderService "github.com/space-wanderer/microservices/order/internal/service/order"
	migrator "github.com/space-wanderer/microservices/platform/pkg/migrator/pg"
	order_v1 "github.com/space-wanderer/microservices/shared/pkg/api/order/v1"
	inventory_v1 "github.com/space-wanderer/microservices/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/space-wanderer/microservices/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API order_v1.Handler

	orderService service.OrderService

	orderRepository repository.OrderRepository

	inventoryClient inventory_v1.InventoryServiceClient
	paymentClient   payment_v1.PaymentServiceClient

	inventoryGRPCClient grpcClient.InventoryClient
	paymentGRPCClient   grpcClient.PaymentClient

	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn

	pgPool *pgxpool.Pool

	pgMigrator *migrator.Migrator
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) order_v1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewAPI(d.OrderService(ctx))
	}
	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewOrderService(d.OrderRepository(ctx), d.InventoryGRPCClient(ctx), d.PaymentGRPCClient(ctx))
	}
	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PGPool(ctx))
	}
	return d.orderRepository
}

func (d *diContainer) InventoryGRPCClient(ctx context.Context) grpcClient.InventoryClient {
	if d.inventoryGRPCClient == nil {
		d.inventoryGRPCClient = grpcClient.NewInventoryClient(d.InventoryClient(ctx))
	}
	return d.inventoryGRPCClient
}

func (d *diContainer) PaymentGRPCClient(ctx context.Context) grpcClient.PaymentClient {
	if d.paymentGRPCClient == nil {
		d.paymentGRPCClient = grpcClient.NewPaymentClient(d.PaymentClient(ctx))
	}
	return d.paymentGRPCClient
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pgPool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			log.Printf("❌ Ошибка подключения к PostgreSQL: %v", err)
			return nil
		}
		d.pgPool = pgPool
	}
	return d.pgPool
}

func (d *diContainer) PGMigrator(ctx context.Context) *migrator.Migrator {
	if d.pgMigrator == nil {
		db := stdlib.OpenDBFromPool(d.PGPool(ctx))
		d.pgMigrator = migrator.NewMigrator(db, config.AppConfig().Postgres.MigrationDir())
	}
	return d.pgMigrator
}

func (d *diContainer) InventoryClient(ctx context.Context) inventory_v1.InventoryServiceClient {
	if d.inventoryClient == nil {
		if d.inventoryConn == nil {
			conn, err := grpc.NewClient(
				config.AppConfig().OrderInventoryGRPC.Address(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				log.Printf("❌ Ошибка подключения к inventory service: %v", err)
				return nil
			}
			d.inventoryConn = conn
		}
		d.inventoryClient = inventory_v1.NewInventoryServiceClient(d.inventoryConn)
	}
	return d.inventoryClient
}

func (d *diContainer) PaymentClient(ctx context.Context) payment_v1.PaymentServiceClient {
	if d.paymentClient == nil {
		if d.paymentConn == nil {
			conn, err := grpc.NewClient(
				config.AppConfig().OrderPaymentGRPC.Address(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				log.Printf("❌ Ошибка подключения к payment service: %v", err)
				return nil
			}
			d.paymentConn = conn
		}
		d.paymentClient = payment_v1.NewPaymentServiceClient(d.paymentConn)
	}
	return d.paymentClient
}
