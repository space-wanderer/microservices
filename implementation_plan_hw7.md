# 📋 Детальный план реализации ДЗ №7: Логирование, метрики и трейсинг

## 🎯 Обзор задачи

Необходимо реализовать продвинутое логирование (Zap + OpenTelemetry Collector → Elasticsearch + Kibana), метрики (Prometheus + Grafana) и распределённый трейсинг (OpenTelemetry) для микросервисной архитектуры.

---

## 📊 Анализ текущего состояния

### ✅ Что уже есть:
- **Платформенный логгер** в `platform/pkg/logger` с поддержкой контекста и trace_id
- **Модули metrics и tracing** в `platform/pkg/` (базовая структура)
- **Сервисы**: order, inventory, payment, assembly, notification, iam
- **Инфраструктура**: Kafka, базовая Docker Compose структура

### ❌ Что нужно добавить:
- Интеграция логгера с OpenTelemetry Collector
- Полная настройка инфраструктуры (OTEL Collector, Elasticsearch, Kibana, Prometheus, Grafana, Jaeger)
- Интеграция метрик в сервисы
- Интеграция трейсинга в сервисы

---

## 🏗️ План реализации

### 🔸 Этап 1: Инфраструктура логов (OTEL Collector + Elasticsearch + Kibana)

#### 1.1 Обновление Docker Compose
**Файл**: `deploy/compose/core/docker-compose.yml`

**Действия**:
- Добавить сервисы:
  - `otel-collector` (образ `otel/opentelemetry-collector-contrib:0.123.0`)
  - `elasticsearch` (образ `elasticsearch:9.15.0`)
  - `kibana` (образ `kibana:9.15.0`)
- Настроить порты и зависимости
- Добавить volumes для данных

#### 1.2 Конфигурация OTEL Collector
**Файл**: `deploy/otel/collector.yaml`

**Структура конфигурации**:
```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    send_batch_size: 1000
    timeout: 10s
  resource:
    attributes:
      - key: service.name
        from_attribute: service.name
        action: insert

exporters:
  elasticsearch:
    endpoint: http://elasticsearch:9200
    index: otel-logs-%{+yyyy.MM.dd}
    mapping:
      mode: ecs

service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch, resource]
      exporters: [elasticsearch]
```

#### 1.3 Переменные окружения
**Файл**: `deploy/env/core.env.template`

**Добавить**:
```bash
# OTEL Collector
OTEL_GRPC_PORT=4317
OTEL_HTTP_PORT=4318
OTEL_HEALTH_CHECK_PORT=13133

# Elasticsearch
ELASTICSEARCH_PORT=9200
ELASTICSEARCH_PASSWORD=changeme

# Kibana
KIBANA_PORT=5601
```

---

### 🔸 Этап 2: Расширение платформенного логгера

#### 2.1 Обновление конфигурации логгера
**Файл**: `platform/pkg/logger/logger.go`

**Изменения**:
- Добавить поддержку множественных выходов (stdout + OTLP)
- Создать OTLP лог-экспортер
- Объединить два zapcore.Core в один логгер
- Добавить корректный shutdown экспортера

**Новая структура**:
```go
type Config interface {
    Outputs() []string // ["stdout", "otlp"]
    OTELCollectorEndpoint() string
    ServiceName() string
    LogLevel() string
}

type logger struct {
    zapLogger *zap.Logger
    exporter  *otlplog.Exporter
}
```

#### 2.2 Создание OTLP лог-экспортера
**Новый файл**: `platform/pkg/logger/otlp_exporter.go`

**Функциональность**:
- Создание OTLP/gRPC экспортера
- Настройка ресурса с метаданными сервиса
- Интеграция с zapcore

#### 2.3 Обновление интерфейсов
**Файл**: `platform/pkg/logger/interfaces.go`

**Добавить**:
```go
type LoggerConfig interface {
    Outputs() []string
    OTELCollectorEndpoint() string
    ServiceName() string
    LogLevel() string
}
```

---

### 🔸 Этап 3: Интеграция нового логгера во все сервисы

#### 3.1 Обновление конфигурации сервисов
**Файлы для обновления**:
- `order/internal/config/env/logger.go`
- `inventory/internal/config/env/logger.go`
- `payment/internal/config/env/logger.go`
- `assembly/internal/config/env/logger.go`
- `notification/internal/config/env/logger.go`
- `iam/internal/config/env/logger.go`

**Добавить поля**:
```go
type loggerEnvConfig struct {
    Level                string `env:"LOG_LEVEL,required"`
    AsJson              bool   `env:"LOGGER_AS_JSON,required"`
    Outputs             string `env:"LOG_OUTPUTS,required"` // "stdout,otlp"
    OTELCollectorEndpoint string `env:"OTEL_COLLECTOR_ENDPOINT,required"`
    ServiceName         string `env:"SERVICE_NAME,required"`
}
```

#### 3.2 Обновление переменных окружения
**Файлы**: `deploy/env/*.env.template`

**Добавить в каждый**:
```bash
# Логирование
LOG_LEVEL=info
LOGGER_AS_JSON=true
LOG_OUTPUTS=stdout,otlp
OTEL_COLLECTOR_ENDPOINT=otel-collector:4317
SERVICE_NAME=order # или inventory, payment, etc.
```

#### 3.3 Обновление инициализации в сервисах
**Файлы**: `*/internal/app/app.go`

**Изменения в initLogger**:
```go
func (a *App) initLogger(ctx context.Context) error {
    return logger.InitWithConfig(config.AppConfig().Logger)
}
```

---

### 🔸 Этап 4: Метрики (Prometheus + Grafana)

#### 4.1 Обновление Docker Compose
**Файл**: `deploy/compose/core/docker-compose.yml`

**Добавить сервисы**:
- `prometheus` (образ `prom/prometheus:v3.3.1`)
- `grafana` (образ `grafana/grafana:12.0.0`)

#### 4.2 Конфигурация Prometheus
**Файл**: `deploy/prometheus/prometheus.yml`

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'order-service'
    static_configs:
      - targets: ['order:8080']
    metrics_path: /metrics
    scrape_interval: 5s

  - job_name: 'assembly-service'
    static_configs:
      - targets: ['assembly:8080']
    metrics_path: /metrics
    scrape_interval: 5s
```

#### 4.3 Интеграция метрик в OrderService
**Файл**: `order/internal/service/order/service.go`

**Добавить метрики**:
```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/metric"
)

type OrderService struct {
    // ... существующие поля
    ordersTotal      metric.Int64Counter
    ordersRevenueTotal metric.Float64Counter
}

func NewOrderService(...) *OrderService {
    meter := otel.Meter("order-service")
    
    ordersTotal, _ := meter.Int64Counter(
        "orders_total",
        metric.WithDescription("Total number of orders created"),
    )
    
    ordersRevenueTotal, _ := meter.Float64Counter(
        "orders_revenue_total",
        metric.WithDescription("Total revenue from orders"),
    )
    
    return &OrderService{
        // ... существующие поля
        ordersTotal:      ordersTotal,
        ordersRevenueTotal: ordersRevenueTotal,
    }
}

func (s *OrderService) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (*order_v1.CreateOrderResponse, error) {
    // ... существующая логика
    
    // Инкрементируем метрики
    s.ordersTotal.Add(ctx, 1)
    s.ordersRevenueTotal.Add(ctx, float64(totalAmount))
    
    return response, nil
}
```

#### 4.4 Интеграция метрик в AssemblyService
**Файл**: `assembly/internal/service/service.go`

**Добавить метрику**:
```go
type AssemblyService struct {
    // ... существующие поля
    assemblyDuration metric.Float64Histogram
}

func NewAssemblyService(...) *AssemblyService {
    meter := otel.Meter("assembly-service")
    
    assemblyDuration, _ := meter.Float64Histogram(
        "assembly_duration_seconds",
        metric.WithDescription("Duration of assembly operations"),
        metric.WithUnit("s"),
    )
    
    return &AssemblyService{
        // ... существующие поля
        assemblyDuration: assemblyDuration,
    }
}

func (s *AssemblyService) ProcessAssembly(ctx context.Context, orderID string) error {
    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        s.assemblyDuration.Record(ctx, duration)
    }()
    
    // ... существующая логика
}
```

#### 4.5 Конфигурация Grafana
**Файл**: `deploy/grafana/provisioning/dashboards/microservices.json`

**Создать дашборд с графиками**:
- `orders_total` - счётчик созданных заказов
- `orders_revenue_total` - суммарная выручка
- `assembly_duration_seconds` - гистограмма длительности сборки

---

### 🔸 Этап 5: Трейсинг (Distributed Tracing)

#### 5.1 Обновление Docker Compose
**Файл**: `deploy/compose/core/docker-compose.yml`

**Добавить сервис**:
- `jaeger` (образ `jaegertracing/jaeger:2.6.0`)

#### 5.2 Обновление конфигурации OTEL Collector
**Файл**: `deploy/otel/collector.yaml`

**Добавить в exporters**:
```yaml
exporters:
  otlp/jaeger:
    endpoint: jaeger:4317
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch, resource]
      exporters: [otlp/jaeger]
```

#### 5.3 Интеграция трейсинга в OrderService
**Файл**: `order/internal/service/order/service.go`

**Добавить трейсинг**:
```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func (s *OrderService) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest) (*order_v1.PayOrderResponse, error) {
    // Создаем root span для операции оплаты
    tracer := otel.Tracer("order-service")
    ctx, span := tracer.Start(ctx, "PayOrder")
    defer span.End()
    
    // Добавляем атрибуты к span
    span.SetAttributes(
        attribute.String("order.id", req.OrderId),
        attribute.Float64("order.amount", float64(req.Amount)),
    )
    
    // Вызываем PaymentService с контекстом трейсинга
    paymentResp, err := s.paymentClient.PayOrder(ctx, &payment_v1.PayOrderRequest{
        OrderId: req.OrderId,
        Amount:  req.Amount,
    })
    
    if err != nil {
        span.RecordError(err)
        return nil, err
    }
    
    span.SetAttributes(
        attribute.String("payment.status", paymentResp.Status),
    )
    
    return &order_v1.PayOrderResponse{
        Success: paymentResp.Success,
    }, nil
}
```

#### 5.4 Интеграция трейсинга в PaymentService
**Файл**: `payment/internal/service/payment/service.go`

**Добавить трейсинг**:
```go
func (s *PaymentService) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
    // Продолжаем трейс из OrderService
    tracer := otel.Tracer("payment-service")
    ctx, span := tracer.Start(ctx, "ProcessPayment")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("payment.order_id", req.OrderId),
        attribute.Float64("payment.amount", float64(req.Amount)),
    )
    
    // ... существующая логика обработки платежа
    
    span.SetAttributes(
        attribute.String("payment.result", "success"),
    )
    
    return &payment_v1.PayOrderResponse{
        Success: true,
        Status:  "completed",
    }, nil
}
```

#### 5.5 Настройка gRPC интерцепторов
**Файл**: `order/internal/app/app.go`

**Добавить в initGRPCClients**:
```go
import "github.com/space-wanderer/microservices/platform/pkg/tracing"

func (a *App) initGRPCClients(ctx context.Context) error {
    // Настройка трейсинга для gRPC клиентов
    paymentConn, err := grpc.Dial(
        config.AppConfig().OrderPaymentGRPC.Address(),
        grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("order-service")),
        grpc.WithInsecure(),
    )
    
    // ... аналогично для других клиентов
}
```

---

### 🔸 Этап 6: Обновление зависимостей

#### 6.1 Обновление go.mod файлов
**Файлы**: `platform/go.mod`, `*/go.mod`

**Добавить зависимости**:
```go
require (
    go.opentelemetry.io/otel v1.32.0
    go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v1.32.0
    go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.32.0
    go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.32.0
    go.opentelemetry.io/otel/sdk v1.32.0
    go.opentelemetry.io/otel/sdk/metric v1.32.0
    go.opentelemetry.io/otel/sdk/trace v1.32.0
)
```

---

### 🔸 Этап 7: Тестирование и валидация

#### 7.1 Проверка логов
- Логи видны в stdout
- Логи попадают в Elasticsearch
- Логи отображаются в Kibana (индекс `otel-logs-*`)

#### 7.2 Проверка метрик
- Метрики доступны в Prometheus
- Графики отображаются в Grafana
- Метрики корректно инкрементируются

#### 7.3 Проверка трейсинга
- Трейсы видны в Jaeger UI
- Цепочка: `HTTP → OrderService → gRPC → PaymentService`
- Контекст корректно передается между сервисами

---

## 📁 Структура файлов после реализации

```
deploy/
├── compose/
│   └── core/
│       └── docker-compose.yml          # Обновлен
├── env/
│   ├── core.env.template               # Обновлен
│   └── *.env.template                  # Обновлены
├── otel/
│   └── collector.yaml                  # Новый
├── prometheus/
│   └── prometheus.yml                  # Новый
└── grafana/
    └── provisioning/
        └── dashboards/
            └── microservices.json      # Новый

platform/
└── pkg/
    ├── logger/
    │   ├── logger.go                   # Обновлен
    │   ├── otlp_exporter.go            # Новый
    │   └── interfaces.go               # Новый
    ├── metrics/
    │   └── metrics.go                  # Уже существует
    └── tracing/
        ├── tracer.go                   # Уже существует
        ├── grpc_interceptor.go         # Уже существует
        └── metadata_carrier.go         # Уже существует

order/
├── internal/
│   ├── config/
│   │   └── env/
│   │       └── logger.go               # Обновлен
│   ├── app/
│   │   └── app.go                      # Обновлен
│   └── service/
│       └── order/
│           └── service.go              # Обновлен
└── go.mod                              # Обновлен

# Аналогично для inventory, payment, assembly, notification, iam
```

---

## 🚀 Порядок выполнения

1. **Инфраструктура** (Этапы 1, 4.1, 5.1) - настройка Docker Compose и конфигураций
2. **Платформенный логгер** (Этап 2) - расширение функциональности
3. **Интеграция логгера** (Этап 3) - обновление всех сервисов
4. **Метрики** (Этапы 4.2-4.5) - интеграция в сервисы и настройка Grafana
5. **Трейсинг** (Этапы 5.2-5.5) - интеграция в сервисы
6. **Зависимости** (Этап 6) - обновление go.mod
7. **Тестирование** (Этап 7) - проверка всех компонентов

---

## ⚠️ Важные замечания

1. **Совместимость API**: Сохранить существующий API логгера для обратной совместимости
2. **Формат логов**: Строго JSON для корректного парсинга в Elasticsearch
3. **Семплирование**: В продакшене настроить семплирование трейсов (например, 10%)
4. **Безопасность**: В продакшене включить TLS для всех соединений
5. **Мониторинг**: Настроить алерты в Grafana для критических метрик

---

## 📚 Полезные ссылки

- [OpenTelemetry Go Documentation](https://opentelemetry.io/docs/languages/go/)
- [Zap Logger Documentation](https://pkg.go.dev/go.uber.org/zap)
- [Prometheus Go Client](https://pkg.go.dev/github.com/prometheus/client_golang)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [Elasticsearch Documentation](https://www.elastic.co/guide/en/elasticsearch/reference/current/)
