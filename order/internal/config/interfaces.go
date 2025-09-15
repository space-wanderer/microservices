package config

type LoggerConfig interface {
	Level() string
	AsJSON() bool
	Outputs() []string
	OTELCollectorEndpoint() string
	ServiceName() string
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

type IamGRPCConfig interface {
	Address() string
}
