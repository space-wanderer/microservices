package env

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type redisEnvConfig struct {
	Host              string        `env:"REDIS_HOST,required"`
	Port              string        `env:"REDIS_PORT,required"`
	ConnectionTimeout time.Duration `env:"REDIS_CONNECTION_TIMEOUT,required"`
	MaxIdle           int           `env:"REDIS_MAX_IDLE,required"`
	IdleTimeout       time.Duration `env:"REDIS_IDLE_TIMEOUT,required"`
}

type RedisConfig struct {
	raw redisEnvConfig
}

func NewRedisConfig() (*RedisConfig, error) {
	var raw redisEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &RedisConfig{raw: raw}, nil
}

func (cfg *RedisConfig) Address() string {
	return fmt.Sprintf("%s:%s", cfg.raw.Host, cfg.raw.Port)
}

func (cfg *RedisConfig) ConnectionTimeout() time.Duration {
	return cfg.raw.ConnectionTimeout
}

func (cfg *RedisConfig) MaxIdle() int {
	return cfg.raw.MaxIdle
}

func (cfg *RedisConfig) IdleTimeout() time.Duration {
	return cfg.raw.IdleTimeout
}
