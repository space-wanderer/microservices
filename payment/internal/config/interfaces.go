package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PaymentConfig interface {
	Address() string
}
