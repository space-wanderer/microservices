package env

import (
	"strings"

	"github.com/caarlos0/env/v11"
)

type loggerEnvConfig struct {
	Level                 string `env:"LOGGER_LEVEL,required"`
	AsJson                bool   `env:"LOGGER_AS_JSON,required"`
	Outputs               string `env:"LOGGER_OUTPUTS" envDefault:"stdout,otlp"`
	OTELCollectorEndpoint string `env:"OTEL_COLLECTOR_ENDPOINT" envDefault:"localhost:4317"`
	ServiceName           string `env:"SERVICE_NAME" envDefault:"payment"`
}

type loggerConfig struct {
	raw loggerEnvConfig
}

func NewLoggerConfig() (*loggerConfig, error) {
	var raw loggerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &loggerConfig{raw: raw}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.raw.Level
}

func (cfg *loggerConfig) AsJson() bool {
	return cfg.raw.AsJson
}

func (cfg *loggerConfig) AsJSON() bool {
	return cfg.raw.AsJson
}

func (cfg *loggerConfig) Outputs() []string {
	// Парсим строку outputs в слайс
	outputs := strings.Split(cfg.raw.Outputs, ",")
	var result []string
	for _, output := range outputs {
		output = strings.TrimSpace(output)
		if output != "" {
			result = append(result, output)
		}
	}
	return result
}

func (cfg *loggerConfig) OTELCollectorEndpoint() string {
	return cfg.raw.OTELCollectorEndpoint
}

func (cfg *loggerConfig) ServiceName() string {
	return cfg.raw.ServiceName
}
