package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamGrpcEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type iamGrpcConfig struct {
	raw iamGrpcEnvConfig
}

func NewIamGrpcConfig() (*iamGrpcConfig, error) {
	var raw iamGrpcEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &iamGrpcConfig{raw: raw}, nil
}

func (cfg *iamGrpcConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
