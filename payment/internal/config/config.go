package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/space-wanderer/microservices/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	PaymentGRPC PaymentConfig
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

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:      loggerCfg,
		PaymentGRPC: paymentGRPCCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
