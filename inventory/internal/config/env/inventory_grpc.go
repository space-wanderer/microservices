package env

import "github.com/caarlos0/env/v11"

type inventoryGRPCEnvConfig struct {
	Port string `env:"GRPC_PORT,required"`
}

type inventoryGRPCConfig struct {
	raw inventoryGRPCEnvConfig
}

func NewInventoryGRPCConfig() (*inventoryGRPCConfig, error) {
	var raw inventoryGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &inventoryGRPCConfig{raw: raw}, nil
}

func (cfg *inventoryGRPCConfig) Address() string {
	return cfg.raw.Port
}
