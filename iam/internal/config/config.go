package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/space-wanderer/microservices/iam/internal/config/env"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
)

var appConfig *config

type config struct {
	Logger   LoggerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Session  SessionConfig
	IamGrpc  IamGrpcConfig
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

	// Инициализируем логгер
	if err := logger.Init(loggerCfg.Level(), loggerCfg.AsJSON()); err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	sessionCfg, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	iamGrpcCfg, err := env.NewIamGrpcConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:   loggerCfg,
		Postgres: postgresCfg,
		Redis:    redisCfg,
		Session:  sessionCfg,
		IamGrpc:  iamGrpcCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
