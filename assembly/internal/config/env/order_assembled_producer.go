package env

import "github.com/caarlos0/env/v11"

type orderAssembledProducerEnvConfig struct {
	TopicName string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
}

type orderAssembledProducerConfig struct {
	raw orderAssembledProducerEnvConfig
}

func NewOrderAssembledProducerConfig() (*orderAssembledProducerConfig, error) {
	var raw orderAssembledProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderAssembledProducerConfig{raw: raw}, nil
}

func (cfg *orderAssembledProducerConfig) TopicName() string {
	return cfg.raw.TopicName
}
