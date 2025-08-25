package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/space-wanderer/microservices/order/internal/api/order/v1"
	grpcClient "github.com/space-wanderer/microservices/order/internal/client/grpc"
	"github.com/space-wanderer/microservices/order/internal/config"
	kafkaConverter "github.com/space-wanderer/microservices/order/internal/converter/kafka"
	orderDecoder "github.com/space-wanderer/microservices/order/internal/converter/kafka/decoder"
	orderProducer "github.com/space-wanderer/microservices/order/internal/converter/kafka/producer"
	"github.com/space-wanderer/microservices/order/internal/repository"
	orderRepository "github.com/space-wanderer/microservices/order/internal/repository/order"
	"github.com/space-wanderer/microservices/order/internal/service"
	orderConsumer "github.com/space-wanderer/microservices/order/internal/service/consumer/order_consumer"
	orderService "github.com/space-wanderer/microservices/order/internal/service/order"
	platformKafka "github.com/space-wanderer/microservices/platform/pkg/kafka"
	"github.com/space-wanderer/microservices/platform/pkg/kafka/consumer"
	"github.com/space-wanderer/microservices/platform/pkg/kafka/producer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
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

	// Kafka Producer для OrderPaidEvent
	orderPaidProducer        platformKafka.Producer
	orderPaidProducerService kafkaConverter.OrderPaidProducer

	// Kafka Consumer для ShipAssembledEvent
	shipAssembledConsumer        platformKafka.Consumer
	shipAssembledDecoder         kafkaConverter.ShipAssembledDecoder
	shipAssembledConsumerService *orderConsumer.Service
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
		d.orderService = orderService.NewOrderService(d.OrderRepository(ctx), d.InventoryGRPCClient(ctx), d.PaymentGRPCClient(ctx), d.OrderPaidProducerService(ctx))
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

// OrderPaidProducer создает Kafka producer для отправки OrderPaidEvent
func (d *diContainer) OrderPaidProducer(ctx context.Context) platformKafka.Producer {
	if d.orderPaidProducer == nil {
		cfg := config.AppConfig()

		// Создаем Sarama producer
		saramaConfig := sarama.NewConfig()
		saramaConfig.Producer.Return.Successes = true
		saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
		saramaConfig.Producer.Retry.Max = 3

		saramaProducer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers(), saramaConfig)
		if err != nil {
			log.Printf("❌ Ошибка создания Sarama producer: %v", err)
			return nil
		}

		// Создаем platform producer
		d.orderPaidProducer = producer.NewProducer(saramaProducer, cfg.OrderPaidProducer.TopicName(), logger.Logger())
	}
	return d.orderPaidProducer
}

// OrderPaidProducerService создает сервис для отправки OrderPaidEvent
func (d *diContainer) OrderPaidProducerService(ctx context.Context) kafkaConverter.OrderPaidProducer {
	if d.orderPaidProducerService == nil {
		d.orderPaidProducerService = orderProducer.NewOrderProducer(d.OrderPaidProducer(ctx))
	}
	return d.orderPaidProducerService
}

// ShipAssembledConsumer создает Kafka consumer для получения ShipAssembledEvent
func (d *diContainer) ShipAssembledConsumer(ctx context.Context) platformKafka.Consumer {
	if d.shipAssembledConsumer == nil {
		cfg := config.AppConfig()

		// Создаем Sarama consumer
		saramaConfig := sarama.NewConfig()
		saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

		saramaConsumer, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers(), cfg.OrderAssembledConsumer.ConsumerGroupID(), saramaConfig)
		if err != nil {
			log.Printf("❌ Ошибка создания Sarama consumer: %v", err)
			return nil
		}

		// Создаем platform consumer
		d.shipAssembledConsumer = consumer.NewConsumer(saramaConsumer, []string{cfg.OrderAssembledConsumer.TopicName()}, logger.Logger())
	}
	return d.shipAssembledConsumer
}

// ShipAssembledDecoder создает decoder для ShipAssembledEvent
func (d *diContainer) ShipAssembledDecoder(ctx context.Context) kafkaConverter.ShipAssembledDecoder {
	if d.shipAssembledDecoder == nil {
		d.shipAssembledDecoder = orderDecoder.NewShipAssembledDecoder()
	}
	return d.shipAssembledDecoder
}

// ShipAssembledConsumerService создает сервис для обработки ShipAssembledEvent
func (d *diContainer) ShipAssembledConsumerService(ctx context.Context) *orderConsumer.Service {
	if d.shipAssembledConsumerService == nil {
		d.shipAssembledConsumerService = orderConsumer.NewService(
			d.ShipAssembledConsumer(ctx),
			d.ShipAssembledDecoder(ctx),
			d.OrderService(ctx),
		)
	}
	return d.shipAssembledConsumerService
}
