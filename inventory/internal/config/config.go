package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/space-wanderer/microservices/inventory/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	InventoryGRPC InventoryGRPCConfig
	Mongo         MongoConfig
	IamGRPC       IamGRPCConfig
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

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}
	iamGRPCCfg, err := env.NewIamGrpcConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		InventoryGRPC: inventoryGRPCCfg,
		Mongo:         mongoCfg,
		IamGRPC:       iamGRPCCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
