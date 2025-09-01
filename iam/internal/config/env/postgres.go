package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type postgresEnvConfig struct {
	Host         string `env:"POSTGRES_HOST,required"`
	Port         string `env:"POSTGRES_PORT,required"`
	Password     string `env:"POSTGRES_PASSWORD,required"`
	Database     string `env:"POSTGRES_DB,required"`
	User         string `env:"POSTGRES_USER,required"`
	MigrationDir string `env:"MIGRATION_DIRECTORY,required"`
}

type PostgresConfig struct {
	raw postgresEnvConfig
}

func NewPostgresConfig() (*PostgresConfig, error) {
	var raw postgresEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &PostgresConfig{raw: raw}, nil
}

func (cfg *PostgresConfig) URI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.raw.User,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.Database,
	)
}

func (cfg *PostgresConfig) Database() string {
	return cfg.raw.Database
}

func (cfg *PostgresConfig) MigrationDir() string {
	return cfg.raw.MigrationDir
}
