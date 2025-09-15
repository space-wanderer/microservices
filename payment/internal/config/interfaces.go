package config

type LoggerConfig interface {
	Level() string
	AsJSON() bool
	Outputs() []string
	OTELCollectorEndpoint() string
	ServiceName() string
}

type PaymentConfig interface {
	Address() string
}
