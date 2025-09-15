package config

type LoggerConfig interface {
	Level() string
	AsJSON() bool
	Outputs() []string
	OTELCollectorEndpoint() string
	ServiceName() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidConsumerConfig interface {
	TopicName() string
	ConsumerGroupID() string
}

type OrderAssembledProducerConfig interface {
	TopicName() string
}
