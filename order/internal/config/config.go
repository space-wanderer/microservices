package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/space-wanderer/microservices/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger             LoggerConfig
	OrderHTTP          OrderHTTPConfig
	OrderPaymentGRPC   OrderPaymentGRPCConfig
	OrderInventoryGRPC OrderInventoryGRPCConfig
	Postgres           PosgresConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	orderHTTPConfig, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	orderPaymentGRPCConfig, err := env.NewOrderPaymentGRPCConfig()
	if err != nil {
		return err
	}

	orderInventoryGRPCConfig, err := env.NewOrderInventoryGRPCConfig()
	if err != nil {
		return err
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:             loggerCfg,
		OrderHTTP:          orderHTTPConfig,
		OrderPaymentGRPC:   orderPaymentGRPCConfig,
		OrderInventoryGRPC: orderInventoryGRPCConfig,
		Postgres:           postgresConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
