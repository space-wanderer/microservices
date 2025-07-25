package env

import (
	"github.com/caarlos0/env/v11"
)

type paymentGRPCEnvConfig struct {
	Port string `env:"GRPC_PORT,required"`
}

type paymentGRPCConfig struct {
	raw paymentGRPCEnvConfig
}

func NewPaymentGRPCConfig() (*paymentGRPCConfig, error) {
	var raw paymentGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &paymentGRPCConfig{raw: raw}, nil
}

func (cfg *paymentGRPCConfig) Address() string {
	return cfg.raw.Port
}
