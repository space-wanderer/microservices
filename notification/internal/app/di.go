package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	"github.com/space-wanderer/microservices/notification/internal/client/http/telegram"
	"github.com/space-wanderer/microservices/notification/internal/config"
	"github.com/space-wanderer/microservices/notification/internal/converter/kafka"
	"github.com/space-wanderer/microservices/notification/internal/converter/kafka/decoder"
	"github.com/space-wanderer/microservices/notification/internal/service"
	consumerAssembledService "github.com/space-wanderer/microservices/notification/internal/service/consumer/order_assembled_consumer"
	consumerPaidService "github.com/space-wanderer/microservices/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/space-wanderer/microservices/notification/internal/service/telegram"
	platformKafka "github.com/space-wanderer/microservices/platform/pkg/kafka"
	"github.com/space-wanderer/microservices/platform/pkg/kafka/consumer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

type diContainer struct {
	orderPaidConsumer      platformKafka.Consumer
	orderAssembledConsumer platformKafka.Consumer

	orderPaidDecoder      kafka.OrderPaidDecoder
	orderAssembledDecoder kafka.ShipAssembledDecoder

	telegramBot     *bot.Bot
	telegramClient  *telegram.Client
	telegramService *telegramService.Service
}

func NewDiContainer() *diContainer {
	return &diContainer{}
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

func (d *diContainer) OrderAssembledConsumer(ctx context.Context) platformKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		cfg := config.AppConfig()

		// Создаем Sarama конфигурацию
		saramaConfig := sarama.NewConfig()
		saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

		// Создаем consumer group
		group, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers(), cfg.OrderAssembledConsumer.ConsumerGroupID(), saramaConfig)
		if err != nil {
			log.Printf("❌ Ошибка создания Kafka consumer group: %v", err)
			return nil
		}

		// Создаем consumer
		d.orderAssembledConsumer = consumer.NewConsumer(group, []string{cfg.OrderAssembledConsumer.TopicName()}, logger.Logger())
	}
	return d.orderAssembledConsumer
}

func (d *diContainer) OrderPaidDecoder(ctx context.Context) kafka.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}
	return d.orderPaidDecoder
}

func (d *diContainer) OrderAssembledDecoder(ctx context.Context) kafka.ShipAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}
	return d.orderAssembledDecoder
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		cfg := config.AppConfig()
		b, err := bot.New(cfg.TelegramBot.Token())
		if err != nil {
			log.Printf("❌ Ошибка создания Telegram бота: %v", err)
			return nil
		}
		d.telegramBot = b
	}
	return d.telegramBot
}

func (d *diContainer) TelegramClient(ctx context.Context) *telegram.Client {
	if d.telegramClient == nil {
		d.telegramClient = telegram.NewClient(d.TelegramBot(ctx))
	}
	return d.telegramClient
}

func (d *diContainer) TelegramService(ctx context.Context) *telegramService.Service {
	if d.telegramService == nil {
		cfg := config.AppConfig()
		d.telegramService = telegramService.NewService(d.TelegramClient(ctx), cfg.TelegramBot)
	}
	return d.telegramService
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) service.ConsumerService {
	return consumerPaidService.NewService(d.OrderPaidConsumer(ctx), d.OrderPaidDecoder(ctx), d.TelegramService(ctx))
}

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) service.ConsumerService {
	return consumerAssembledService.NewService(d.OrderAssembledConsumer(ctx), d.OrderAssembledDecoder(ctx), d.TelegramService(ctx))
}
