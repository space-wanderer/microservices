package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/space-wanderer/microservices/assembly/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
	OrderAssembledProducer OrderAssembledProducerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	orderAssembledProducerCfg, err := env.NewOrderAssembledProducerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		Kafka:                  kafkaCfg,
		OrderPaidConsumer:      orderPaidConsumerCfg,
		OrderAssembledProducer: orderAssembledProducerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
