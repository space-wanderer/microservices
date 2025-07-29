package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type orderPaymentGRPCEnvConfig struct {
	Host string `env:"PAYMENT_GRPC_HOST,required"`
	Port string `env:"PAYMENT_GRPC_PORT,required"`
}

type orderPaymentGRPCConfig struct {
	raw orderPaymentGRPCEnvConfig
}

func NewOrderPaymentGRPCConfig() (*orderPaymentGRPCConfig, error) {
	var raw orderPaymentGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderPaymentGRPCConfig{raw: raw}, nil
}

func (cfg *orderPaymentGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
