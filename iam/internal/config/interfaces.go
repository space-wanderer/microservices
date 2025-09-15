package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJSON() bool
	Outputs() []string
	OTELCollectorEndpoint() string
	ServiceName() string
}

type PostgresConfig interface {
	URI() string
	Database() string
	MigrationDir() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type SessionConfig interface {
	SessionTTL() time.Duration
}

type IamGrpcConfig interface {
	Address() string
}
