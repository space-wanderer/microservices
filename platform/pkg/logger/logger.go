// Package logger предоставляет dual-write логгер с использованием zapcore.Tee архитектуры
//
// АРХИТЕКТУРА ЛОГГЕРА:
//
// Логгер использует zapcore.NewTee для параллельной записи в два назначения:
// 1. Stdout (для Kubernetes/контейнерных окружений)
// 2. OpenTelemetry коллектор (для централизованного сбора логов)
//
// ПОТОК ДАННЫХ:
//
//		Application
//		    ↓ (logger.Info/Error)
//		zap.Logger
//		    ↓
//		zapcore.Tee
//		   ↙        ↘
//	 StdoutCore   SimpleOTLPCore
//		   ↓             ↓
//	 os.Stdout   SimpleOTLPWriter
//		               ↓
//		        zapcore.BufferedWriteSyncer
//		               ↓
//		         OTLP Collector (gRPC)
//
// КОМПОНЕНТЫ:
//
// 1. StdoutCore - стандартный zap core для вывода в консоль
// 2. SimpleOTLPCore - преобразует zap Entry в OpenTelemetry Record
// 3. SimpleOTLPWriter - отправляет OTLP Records в коллектор
// 4. BufferedWriteSyncer - буферизация для асинхронной отправки
//
// ОСОБЕННОСТИ:
//
// - Graceful degradation: при недоступности OTLP коллектора stdout продолжает работать
// - Метрики: отслеживание sent/dropped записей для мониторинга
// - Батчирование: OTLP SDK автоматически группирует записи для эффективной отправки
// - Таймауты: 500ms лимит для предотвращения блокировки приложения
package logger

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Глобальные переменные пакета
var (
	global   *zap.Logger     // глобальный экземпляр логгера
	initOnce sync.Once       // обеспечивает единократную инициализацию
	level    zap.AtomicLevel // уровень логирования (может изменяться динамически)
	otlpCore *SimpleOTLPCore // OTLP core для graceful shutdown
)

// Константы конфигурации OTLP (значения по умолчанию)
const (
	defaultOTLPEndpoint       = "localhost:4317" // адрес OTLP коллектора по умолчанию
	defaultServiceName        = "microservice"   // имя сервиса по умолчанию
	defaultServiceEnvironment = "dev"            // окружение для фильтрации логов
)

// Init инициализирует глобальный логгер с Tee архитектурой.
// Поддерживает одновременную запись в stdout и OTLP коллектор.
//
// Параметры:
//   - logLevel: уровень логирования ("debug", "info", "warn", "error")
//   - asJSON: формат вывода (true - JSON, false - консольный)
//   - enableOTLP: включение отправки в OpenTelemetry коллектор
func Init(logLevel string, asJSON bool, enableOTLP bool) error {
	initOnce.Do(func() {
		level = zap.NewAtomicLevelAt(parseLevel(logLevel))
		cores := buildCores(asJSON, enableOTLP)
		global = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
	})

	if global == nil {
		return fmt.Errorf("logger init failed")
	}

	return nil
}

// InitWithConfig инициализирует глобальный логгер с использованием конфигурации.
// Поддерживает множественные выходы согласно конфигурации.
func InitWithConfig(cfg LoggerConfig) error {
	initOnce.Do(func() {
		level = zap.NewAtomicLevelAt(parseLevel(cfg.Level()))
		cores := buildCoresFromConfig(cfg)
		global = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
	})

	if global == nil {
		return fmt.Errorf("logger init failed")
	}

	return nil
}

// buildCores создает слайс cores для zapcore.Tee.
// Всегда включает stdout core, опционально добавляет OTLP core.
func buildCores(asJSON bool, enableOTLP bool) []zapcore.Core {
	cores := []zapcore.Core{
		createStdoutCore(asJSON),
	}

	if enableOTLP {
		if otlpCore := createOTLPCore(); otlpCore != nil {
			cores = append(cores, otlpCore)
		}
	}

	return cores
}

// buildCoresFromConfig создает слайс cores на основе конфигурации.
// Поддерживает множественные выходы согласно настройкам.
func buildCoresFromConfig(cfg LoggerConfig) []zapcore.Core {
	var cores []zapcore.Core

	outputs := cfg.Outputs()
	for _, output := range outputs {
		switch output {
		case "stdout":
			cores = append(cores, createStdoutCore(cfg.AsJSON()))
		case "otlp":
			if otlpCore := createOTLPCoreFromConfig(cfg); otlpCore != nil {
				cores = append(cores, otlpCore)
			}
		}
	}

	// Если выходы не указаны, используем stdout по умолчанию
	if len(cores) == 0 {
		cores = append(cores, createStdoutCore(cfg.AsJSON()))
	}

	return cores
}

// createStdoutCore создает core для записи в stdout/stderr.
// Поддерживает JSON и консольный формат вывода.
func createStdoutCore(asJSON bool) zapcore.Core {
	config := buildEncoderConfig()
	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
}

// createOTLPCore создает core для отправки в OpenTelemetry коллектор.
// При ошибке подключения возвращает nil (graceful degradation).
func createOTLPCore() *SimpleOTLPCore {
	// Создаем упрощенный OTLP core
	core := NewSimpleOTLPCore(defaultOTLPEndpoint, defaultServiceName, level)
	otlpCore = core // Сохраняем ссылку для graceful shutdown
	return core
}

// createOTLPCoreFromConfig создает core для отправки в OpenTelemetry коллектор на основе конфигурации.
// При ошибке подключения возвращает nil (graceful degradation).
func createOTLPCoreFromConfig(cfg LoggerConfig) *SimpleOTLPCore {
	endpoint := cfg.OTELCollectorEndpoint()
	if endpoint == "" {
		endpoint = defaultOTLPEndpoint
	}

	svcName := cfg.ServiceName()
	if svcName == "" {
		svcName = defaultServiceName
	}

	// Создаем упрощенный OTLP core
	core := NewSimpleOTLPCore(endpoint, svcName, level)
	otlpCore = core // Сохраняем ссылку для graceful shutdown
	return core
}

// buildEncoderConfig настраивает формат вывода логов с нужными полями
func buildEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
}

// Info записывает лог уровня INFO.
// Отправляется одновременно в stdout и OTLP коллектор (если включен).
func Info(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Info(msg, fields...)
	}
}

// Error записывает лог уровня ERROR.
// Отправляется одновременно в stdout и OTLP коллектор (если включен).
func Error(_ context.Context, msg string, fields ...zap.Field) {
	if global != nil {
		global.Error(msg, fields...)
	}
}

// Sync принудительно сбрасывает все буферизованные логи.
// Вызывает sync для всех cores (stdout + OTLP).
func Sync() error {
	if global != nil {
		return global.Sync()
	}

	return nil
}

// Close корректно завершает работу логгера.
// Закрывает gRPC соединение и синхронизирует буферы.
func Close() error {
	// Сначала синхронизируем буферы
	if err := Sync(); err != nil {
		return err
	}

	// Затем закрываем OTLP core
	if otlpCore != nil {
		return otlpCore.Close()
	}

	return nil
}

// parseLevel преобразует строковое значение в zapcore.Level
func parseLevel(levelStr string) zapcore.Level {
	switch levelStr {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
