package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OrderHTTPConfig interface {
	Address() string
}

type OrderPaymentGRPCConfig interface {
	Address() string
}

type OrderInventoryGRPCConfig interface {
	Address() string
}

type PosgresConfig interface {
	URI() string
	Database() string
	MigrationDir() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderAssembledConsumerConfig interface {
	TopicName() string
	ConsumerGroupID() string
}

type OrderPaidProducerConfig interface {
	TopicName() string
}
