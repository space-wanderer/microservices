package logger

// LoggerConfig определяет интерфейс конфигурации логгера
type LoggerConfig interface {
	// Outputs возвращает список выходов для логирования
	// Поддерживаемые значения: "stdout", "otlp"
	Outputs() []string

	// OTELCollectorEndpoint возвращает адрес OTLP коллектора
	// Например: "otel-collector:4317"
	OTELCollectorEndpoint() string

	// ServiceName возвращает имя сервиса для телеметрии
	// Например: "order-service", "inventory-service"
	ServiceName() string

	// LogLevel возвращает уровень логирования
	// Поддерживаемые значения: "debug", "info", "warn", "error"
	LogLevel() string

	// AsJSON возвращает true, если логи должны выводиться в JSON формате
	AsJSON() bool
}
