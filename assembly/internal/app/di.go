package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/space-wanderer/microservices/assembly/internal/config"
	"github.com/space-wanderer/microservices/assembly/internal/converter/kafka"
	"github.com/space-wanderer/microservices/assembly/internal/converter/kafka/decoder"
	"github.com/space-wanderer/microservices/assembly/internal/service"
	consumerService "github.com/space-wanderer/microservices/assembly/internal/service/consumer/order_consumer"
	producerService "github.com/space-wanderer/microservices/assembly/internal/service/producer/order_producer"
	platformKafka "github.com/space-wanderer/microservices/platform/pkg/kafka"
	"github.com/space-wanderer/microservices/platform/pkg/kafka/consumer"
	"github.com/space-wanderer/microservices/platform/pkg/kafka/producer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

type diContainer struct {
	consumerService service.ConsumerService
	producerService service.ProducerService

	orderPaidConsumer      platformKafka.Consumer
	orderAssembledProducer platformKafka.Producer

	orderPaidDecoder kafka.AssemblyRecodedDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) ConsumerService(ctx context.Context) service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = consumerService.NewService(d.OrderPaidConsumer(ctx), d.OrderPaidDecoder(ctx))
	}
	return d.consumerService
}

func (d *diContainer) ProducerService(ctx context.Context) service.ProducerService {
	if d.producerService == nil {
		d.producerService = producerService.NewService(d.OrderAssembledProducer(ctx))
	}
	return d.producerService
}

func (d *diContainer) OrderPaidConsumer(ctx context.Context) platformKafka.Consumer {
	if d.orderPaidConsumer == nil {
		cfg := config.AppConfig()

		// Создаем Sarama конфигурацию
		saramaConfig := sarama.NewConfig()
		saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

		// Создаем consumer group
		group, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers(), cfg.OrderPaidConsumer.ConsumerGroupID(), saramaConfig)
		if err != nil {
			log.Printf("❌ Ошибка создания Kafka consumer group: %v", err)
			return nil
		}

		// Создаем consumer
		d.orderPaidConsumer = consumer.NewConsumer(group, []string{cfg.OrderPaidConsumer.TopicName()}, logger.Logger())
	}
	return d.orderPaidConsumer
}

func (d *diContainer) OrderAssembledProducer(ctx context.Context) platformKafka.Producer {
	if d.orderAssembledProducer == nil {
		cfg := config.AppConfig()

		// Создаем Sarama конфигурацию
		saramaConfig := sarama.NewConfig()
		saramaConfig.Producer.Return.Successes = true
		saramaConfig.Producer.RequiredAcks = sarama.WaitForAll

		// Создаем sync producer
		syncProducer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers(), saramaConfig)
		if err != nil {
			log.Printf("❌ Ошибка создания Kafka producer: %v", err)
			return nil
		}

		// Создаем producer
		d.orderAssembledProducer = producer.NewProducer(syncProducer, cfg.OrderAssembledProducer.TopicName(), logger.Logger())
	}
	return d.orderAssembledProducer
}

func (d *diContainer) OrderPaidDecoder(ctx context.Context) kafka.AssemblyRecodedDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}
	return d.orderPaidDecoder
}
