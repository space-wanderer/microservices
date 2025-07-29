package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type orderInventoryGRPCEnvConfig struct {
	Host string `env:"INVENTORY_GRPC_HOST,required"`
	Port string `env:"INVENTORY_GRPC_PORT,required"`
}

type orderInventoryGRPCConfig struct {
	raw orderInventoryGRPCEnvConfig
}

func NewOrderInventoryGRPCConfig() (*orderInventoryGRPCConfig, error) {
	var raw orderInventoryGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderInventoryGRPCConfig{raw: raw}, nil
}

func (cfg *orderInventoryGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
