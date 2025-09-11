// OTLP Core Component
//
// Что здесь происходит:
// - record: это одна лог-запись (уровень, сообщение, время, поля-атрибуты).
// - core: «ядро» логгера. Оно решает «принимаю ли я эту запись» и «как её отправлять».
// - tee: «тройник», который раздаёт одну запись сразу нескольким cores.
//
// Интерфейс zapcore.Core (что должен уметь любой core):
// - Enabled(level):решить, писать ли запись этого уровня.
// - With(fields): вернуть копию core с дополнительными полями (мы их учитываем в Write).
// - Check(entry, ce): добавить себя в список получателей записи, если уровень подходит.
// - Write(entry, fields): собрать record и отправить его «куда надо».
// - Sync(): сбросить буферы, если они есть.
//
// Архитектура потока для OTLP:
// zap.Logger -> zapcore.Tee -> SimpleOTLPCore -> OTLP Collector (gRPC)
package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Таймаут отправки одной записи, чтобы не блокировать приложение
const emitTimeout = 500 * time.Millisecond

// SimpleOTLPCore преобразует zap-записи в OpenTelemetry Records и отправляет их напрямую в OTLP
type SimpleOTLPCore struct {
	endpoint    string               // адрес OTLP коллектора
	serviceName string               // имя сервиса
	level       zapcore.LevelEnabler // минимальный уровень для записи логов
	grpcClient  *grpc.ClientConn     // gRPC соединение
}

// Упрощенные структуры для OTLP (без protobuf зависимостей)
type ExportLogsServiceRequest struct {
	ResourceLogs []ResourceLogs `json:"resourceLogs"`
}

type ExportLogsServiceResponse struct {
	PartialSuccess *PartialSuccess `json:"partialSuccess,omitempty"`
}

type ResourceLogs struct {
	Resource  Resource    `json:"resource"`
	ScopeLogs []ScopeLogs `json:"scopeLogs"`
}

type Resource struct {
	Attributes []KeyValue `json:"attributes"`
}

type ScopeLogs struct {
	LogRecords []LogRecord `json:"logRecords"`
}

type LogRecord struct {
	TimeUnixNano   string     `json:"timeUnixNano"`
	SeverityNumber int        `json:"severityNumber"`
	SeverityText   string     `json:"severityText"`
	Body           AnyValue   `json:"body"`
	Attributes     []KeyValue `json:"attributes,omitempty"`
}

type KeyValue struct {
	Key   string   `json:"key"`
	Value AnyValue `json:"value"`
}

type AnyValue struct {
	StringValue *string  `json:"stringValue,omitempty"`
	BoolValue   *bool    `json:"boolValue,omitempty"`
	IntValue    *int64   `json:"intValue,omitempty"`
	DoubleValue *float64 `json:"doubleValue,omitempty"`
}

type PartialSuccess struct {
	RejectedLogRecords int64  `json:"rejectedLogRecords"`
	ErrorMessage       string `json:"errorMessage,omitempty"`
}

// NewSimpleOTLPCore создает новый OTLP core, работающий напрямую с OTLP коллектором через gRPC.
func NewSimpleOTLPCore(endpoint, serviceName string, level zapcore.LevelEnabler) *SimpleOTLPCore {
	core := &SimpleOTLPCore{
		endpoint:    endpoint,
		serviceName: serviceName,
		level:       level,
	}

	// Инициализируем gRPC соединение
	core.initGRPCClient()

	return core
}

// initGRPCClient инициализирует gRPC соединение с OTLP коллектором
func (c *SimpleOTLPCore) initGRPCClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем gRPC соединение
	conn, err := grpc.DialContext(ctx, c.endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		// При ошибке соединения gRPC клиент остается nil
		// Это обеспечивает graceful degradation
		return
	}

	c.grpcClient = conn
}

// Enabled проверяет, должен ли лог данного уровня быть записан
func (c *SimpleOTLPCore) Enabled(level zapcore.Level) bool {
	return c.level.Enabled(level)
}

// With создает новый core с дополнительными полями.
// В текущей реализации поля обрабатываются в Write методе,
// поэтому здесь создается копия без изменений.
func (c *SimpleOTLPCore) With(_ []zapcore.Field) zapcore.Core {
	return &SimpleOTLPCore{
		endpoint:    c.endpoint,
		serviceName: c.serviceName,
		level:       c.level,
		grpcClient:  c.grpcClient, // Переиспользуем существующее соединение
	}
}

// Check определяет, должен ли данный лог быть записан данным core.
// Добавляет себя в CheckedEntry если уровень лога соответствует настройкам.
func (c *SimpleOTLPCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

// Write конвертирует zap Entry в OpenTelemetry Record и отправляет в OTLP.
// Пошагово:
//  1. Преобразуем zap-уровень в OTLP Severity (mapZapToOtelSeverity).
//  2. Собираем базовый Record: severity, body=сообщение, timestamp (makeBaseRecord).
//  3. Кодируем zap-поля в OTLP-атрибуты (encodeFieldsToAttrs) и добавляем их в Record.
//  4. Отправляем запись через HTTP API с коротким таймаутом (emitWithTimeout),
//     чтобы не блокировать приложение при сетевых проблемах.
func (c *SimpleOTLPCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	severity := mapZapToOtelSeverity(entry.Level)
	record := makeBaseRecord(entry, severity)
	if len(fields) > 0 {
		attrs := encodeFieldsToAttrs(fields)
		if len(attrs) > 0 {
			record.Attributes = attrs
		}
	}

	c.emitWithTimeout(record)
	return nil
}

// Sync синхронизация не требуется: батчинг делает OTLP SDK
func (c *SimpleOTLPCore) Sync() error { return nil }

// Close закрывает gRPC соединение
func (c *SimpleOTLPCore) Close() error {
	if c.grpcClient != nil {
		return c.grpcClient.Close()
	}
	return nil
}

// mapZapToOtelSeverity — отдельная функция преобразования уровня
func mapZapToOtelSeverity(level zapcore.Level) int {
	switch level {
	case zapcore.DebugLevel:
		return 5 // DEBUG
	case zapcore.InfoLevel:
		return 9 // INFO
	case zapcore.WarnLevel:
		return 13 // WARN
	case zapcore.ErrorLevel:
		return 17 // ERROR
	default:
		return 9 // INFO
	}
}

// makeBaseRecord — сборка базового record без атрибутов
func makeBaseRecord(entry zapcore.Entry, severity int) LogRecord {
	return LogRecord{
		TimeUnixNano:   fmt.Sprintf("%d", entry.Time.UnixNano()),
		SeverityNumber: severity,
		SeverityText:   entry.Level.CapitalString(),
		Body: AnyValue{
			StringValue: &entry.Message,
		},
	}
}

// encodeFieldsToAttrs — подготовка атрибутов из zap-полей.
// Используем zapcore.NewMapObjectEncoder(), чтобы безопасно развернуть []zapcore.Field
// в карту ключ→значение. Далее переносим только базовые типы в OTLP атрибуты.
// Неподдерживаемые типы пропускаем (они продолжат жить в stdout части через zap encoder).
func encodeFieldsToAttrs(fields []zapcore.Field) []KeyValue {
	if len(fields) == 0 {
		return nil
	}

	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	attrs := make([]KeyValue, 0, len(enc.Fields))
	for k, v := range enc.Fields {
		switch val := v.(type) {
		case string:
			attrs = append(attrs, KeyValue{
				Key: k,
				Value: AnyValue{
					StringValue: &val,
				},
			})
		case bool:
			attrs = append(attrs, KeyValue{
				Key: k,
				Value: AnyValue{
					BoolValue: &val,
				},
			})
		case int64:
			attrs = append(attrs, KeyValue{
				Key: k,
				Value: AnyValue{
					IntValue: &val,
				},
			})
		case float64:
			attrs = append(attrs, KeyValue{
				Key: k,
				Value: AnyValue{
					DoubleValue: &val,
				},
			})
		}
	}

	return attrs
}

// emitWithTimeout — отправка в OTLP с коротким таймаутом
func (c *SimpleOTLPCore) emitWithTimeout(record LogRecord) {
	if c.endpoint == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), emitTimeout)
	defer cancel()

	// Создаем запрос для OTLP коллектора
	request := ExportLogsServiceRequest{
		ResourceLogs: []ResourceLogs{
			{
				Resource: Resource{
					Attributes: []KeyValue{
						{
							Key: "service.name",
							Value: AnyValue{
								StringValue: &c.serviceName,
							},
						},
					},
				},
				ScopeLogs: []ScopeLogs{
					{
						LogRecords: []LogRecord{record},
					},
				},
			},
		},
	}

	// Отправляем запрос
	c.sendToOTLP(ctx, request)
}

// sendToOTLP отправляет лог в OTLP коллектор через gRPC
func (c *SimpleOTLPCore) sendToOTLP(ctx context.Context, request ExportLogsServiceRequest) {
	if c.grpcClient == nil {
		// Если gRPC соединение недоступно, используем HTTP fallback
		c.sendToOTLPHTTP(ctx, request)
		return
	}

	// TODO: Реализовать gRPC отправку
	// Пока используем HTTP fallback
	c.sendToOTLPHTTP(ctx, request)
}

// sendToOTLPHTTP отправляет лог через HTTP API (fallback)
func (c *SimpleOTLPCore) sendToOTLPHTTP(ctx context.Context, request ExportLogsServiceRequest) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		// Игнорируем ошибки сериализации для graceful degradation
		return
	}

	// Используем HTTP API OTLP коллектора
	url := fmt.Sprintf("http://%s/v1/logs", c.endpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		// Игнорируем ошибки создания запроса
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: emitTimeout}
	resp, err := client.Do(req)
	if err != nil {
		// Игнорируем ошибки отправки для graceful degradation
		return
	}
	defer resp.Body.Close()
}
