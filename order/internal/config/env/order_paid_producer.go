package env

import "github.com/caarlos0/env/v11"

type orderPaidProducerEnvConfig struct {
	TopicName string `env:"ORDER_PAID_TOPIC_NAME,required"`
}

type orderPaidProducerConfig struct {
	raw orderPaidProducerEnvConfig
}

func NewOrderPaidProducerConfig() (*orderPaidProducerConfig, error) {
	var raw orderPaidProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaidProducerConfig{raw: raw}, nil
}

func (cfg *orderPaidProducerConfig) TopicName() string {
	return cfg.raw.TopicName
}
