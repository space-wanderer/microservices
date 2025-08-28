# 🧱 Домашнее задание — Неделя 6
**Курс: «Микросервисы, как в BigTech 2.0»**

---

## 📌 Задание

На этой неделе вы создадите новый сервис — `IAMService`, который отвечает за регистрацию, аутентификацию и хранение сессий пользователей. Необходимо реализовать gRPC API, настроить работу с PostgreSQL и Redis и интегрировать миграции.

---

### ✅ Что нужно сделать:

1. Создать новый сервис `IAMService`:
    - Использовать ту же архитектуру, что и в других сервисах: `cmd/`, `internal/`, `config`, `service`, `repository`, `di`
    - Поддержка конфигурации через `.env`
    - Самописный DI-контейнер
2. Реализовать **gRPC API** в соответствии с контрактом в [`contracts/iam_service_contracts.md`](contracts/iam_service_contracts.md)
3. Реализовать хранение данных:
    - Пользователи — в **PostgreSQL**
        - Реализовать миграции SQL
        - Обеспечить автоматическую накатку миграций при старте сервиса
    - Сессии — в **Redis**
        - Срок жизни — **24 часа (TTL)**
4. Подключить аутентификационный grpc интерцептор к `Inventory сервису`
5. Подключить аутентификационную мидлварю к `Order сервису`
6. Не забыть вызвать в клиенте `Inventory сервиса` внутри `Order сервиса` функцию `ForwardSessionUUIDToGRPC`, для обогащения outgoing контекста идентификатором сессии

---

📂 Для упрощения выполнения домашнего задания, в папке [`/boilerplates`](boilerplates) находятся вспомогательные файлы:
- `Taskfile.yml` — файл с готовыми командами для управления проектом (907 строк):
  - Команды для работы с Docker Compose и базами данных
  - Команды для генерации кода и тестирования
  - Команды для управления миграциями
- Папка `cache/` — готовые компоненты для работы с Redis:
  - `client.go` — интерфейс для кеша
  - `redis/` — реализация Redis клиента с операторами
- Папка `middleware/` — middleware для авторизации:
  - `grpc/auth.go` — gRPC interceptor для проверки сессий через metadata
  - `http/auth.go` — HTTP middleware для проверки сессий через заголовки
- Папка `deploy/` — конфигурации для развертывания

---

📁 Контракт для `IAMService` находится в файле [`contracts/iam_service_contracts.md`](contracts/iam_service_contracts.md) и описывает gRPC API для управления пользователями и сессиями: методы `Register`, `Login`, `Whoami`, `GetUser`.

📌 Контракт описывает **структуру gRPC API** — к финалу ваша реализация должна соответствовать этому контракту.

---

## 🛠 Актуальная структура проекта

```
.
├── README.md
├── Taskfile.yml
├── assembly
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   └── internal
│       ├── app
│       │   ├── app.go
│       │   └── di.go
│       ├── config
│       │   ├── config.go
│       │   ├── env
│       │   │   ├── kafka.go
│       │   │   ├── logger.go
│       │   │   ├── order_assembled_producer.go
│       │   │   └── order_paid_consumer.go
│       │   ├── interfaces.go
│       │   └── mocks
│       │       ├── mock_kafka_config.go
│       │       ├── mock_logger_config.go
│       │       ├── mock_order_assembled_producer_config.go
│       │       └── mock_order_paid_consumer_config.go
│       ├── converter
│       │   └── kafka
│       │       ├── decoder
│       │       │   └── order_paid.go
│       │       └── kafka.go
│       ├── model
│       │   └── events.go
│       └── service
│           ├── consumer
│           │   └── order_consumer
│           │       ├── consumer.go
│           │       └── handler.go
│           ├── mocks
│           │   ├── mock_consumer_service.go
│           │   └── mock_order_producer_service.go
│           ├── producer
│           │   └── order_producer
│           │       └── producer.go
│           └── service.go
├── buf.work.yaml
├── deploy
│   ├── compose
│   │   ├── assembly
│   │   ├── core
│   │   │   └── docker-compose.yml
│   │   ├── iam
│   │   │   └── docker-compose.yml
│   │   ├── inventory
│   │   │   └── docker-compose.yml
│   │   ├── notification
│   │   ├── order
│   │   │   └── docker-compose.yml
│   │   └── payment
│   ├── docker
│   │   └── inventory
│   │       └── Dockerfile
│   └── env
│       ├── assembly.env.template
│       ├── core.env.template
│       ├── generate-env.sh
│       ├── iam.env.template
│       ├── inventory.env.template
│       ├── notification.env.template
│       ├── order.env.template
│       └── payment.env.template
├── go.work
├── go.work.sum
├── iam
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── api
│   │   │   ├── auth
│   │   │   │   └── v1
│   │   │   │       ├── api.go
│   │   │   │       ├── login.go
│   │   │   │       └── whoami.go
│   │   │   └── user
│   │   │       └── v1
│   │   │           ├── api.go
│   │   │           ├── get.go
│   │   │           └── register.go
│   │   ├── app
│   │   │   ├── app.go
│   │   │   └── di.go
│   │   ├── config
│   │   │   ├── config.go
│   │   │   ├── env
│   │   │   │   ├── iam_grpc.go
│   │   │   │   ├── logger.go
│   │   │   │   ├── postgres.go
│   │   │   │   ├── redis.go
│   │   │   │   └── session.go
│   │   │   ├── interfaces.go
│   │   │   └── mocks
│   │   │       ├── mock_iam_grpc_config.go
│   │   │       ├── mock_logger_config.go
│   │   │       ├── mock_postgres_config.go
│   │   │       ├── mock_redis_config.go
│   │   │       └── mock_session_config.go
│   │   ├── converter
│   │   │   ├── auth.go
│   │   │   └── user.go
│   │   ├── model
│   │   │   ├── auth.go
│   │   │   ├── errors.go
│   │   │   └── user.go
│   │   ├── repository
│   │   │   ├── converter
│   │   │   │   ├── session.go
│   │   │   │   └── user.go
│   │   │   ├── mocks
│   │   │   │   ├── mock_session_repository.go
│   │   │   │   └── mock_user_repository.go
│   │   │   ├── model
│   │   │   │   ├── session.go
│   │   │   │   └── user.go
│   │   │   ├── repository.go
│   │   │   ├── session
│   │   │   │   ├── add_session_to_user_set.go
│   │   │   │   ├── create.go
│   │   │   │   ├── get.go
│   │   │   │   └── repository.go
│   │   │   └── user
│   │   │       ├── create.go
│   │   │       ├── get.go
│   │   │       └── repository.go
│   │   └── service
│   │       ├── auth
│   │       │   ├── login.go
│   │       │   ├── service.go
│   │       │   └── whoami.go
│   │       ├── mocks
│   │       │   ├── mock_auth_service.go
│   │       │   └── mock_user_service.go
│   │       ├── service.go
│   │       └── user
│   │           ├── get.go
│   │           ├── register.go
│   │           └── service.go
│   └── migrations
│       ├── 20250404191615_create_uuid_ossp_extension.sql
│       └── 20250404191624_create_user_table.sql
├── inventory
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── api
│   │   │   └── inventory
│   │   │       └── v1
│   │   │           ├── api.go
│   │   │           ├── get.go
│   │   │           └── list.go
│   │   ├── app
│   │   │   ├── app.go
│   │   │   └── di.go
│   │   ├── client
│   │   │   └── grpc
│   │   │       ├── client.go
│   │   │       └── iam
│   │   │           └── v1
│   │   │               └── client.go
│   │   ├── config
│   │   │   ├── config.go
│   │   │   ├── env
│   │   │   │   ├── iam_grpc.go
│   │   │   │   ├── inventory_grpc.go
│   │   │   │   ├── logger.go
│   │   │   │   └── mongo.go
│   │   │   ├── interfaces.go
│   │   │   └── mocks
│   │   │       ├── mock_iamgrpc_config.go
│   │   │       ├── mock_inventory_grpc_config.go
│   │   │       ├── mock_logger_config.go
│   │   │       └── mock_mongo_config.go
│   │   ├── converter
│   │   │   ├── part.go
│   │   │   └── user.go
│   │   ├── model
│   │   │   ├── auth.go
│   │   │   ├── const.go
│   │   │   ├── errors.go
│   │   │   ├── part.go
│   │   │   └── user.go
│   │   ├── repository
│   │   │   ├── converter
│   │   │   │   └── part.go
│   │   │   ├── mocks
│   │   │   │   └── mock_part_repository.go
│   │   │   ├── model
│   │   │   │   └── part.go
│   │   │   ├── part
│   │   │   │   ├── get.go
│   │   │   │   ├── init.go
│   │   │   │   ├── list.go
│   │   │   │   └── repository.go
│   │   │   └── repository.go
│   │   └── service
│   │       ├── mocks
│   │       │   └── mock_part_service.go
│   │       ├── part
│   │       │   ├── get.go
│   │       │   ├── get_test.go
│   │       │   ├── list.go
│   │       │   ├── list_test.go
│   │       │   ├── service.go
│   │       │   └── suite_test.go
│   │       └── service.go
│   └── tests
│       └── integration
│           ├── constants.go
│           ├── inventory_test.go
│           ├── setup.go
│           ├── suite_test.go
│           ├── teardown.go
│           └── test_environment.go
├── notification
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   └── internal
│       ├── app
│       │   ├── app.go
│       │   └── di.go
│       ├── client
│       │   └── http
│       │       ├── client.go
│       │       └── telegram
│       │           └── client.go
│       ├── config
│       │   ├── config.go
│       │   ├── env
│       │   │   ├── kafka.go
│       │   │   ├── logger.go
│       │   │   ├── order_assembled_consumer.go
│       │   │   ├── order_paid_consumer.go
│       │   │   └── telegram_bot.go
│       │   ├── interfaces.go
│       │   └── mocks
│       │       ├── mock_kafka_config.go
│       │       ├── mock_logger_config.go
│       │       ├── mock_order_assembled_consumer_config.go
│       │       ├── mock_order_paid_consumer_config.go
│       │       └── mock_telegram_bot_config.go
│       ├── converter
│       │   └── kafka
│       │       ├── decoder
│       │       │   ├── order_assembled.go
│       │       │   └── order_paid.go
│       │       └── kafka.go
│       ├── model
│       │   └── events.go
│       └── service
│           ├── consumer
│           │   ├── order_assembled_consumer
│           │   │   ├── consumer.go
│           │   │   └── handler.go
│           │   └── order_paid_consumer
│           │       ├── consumer.go
│           │       └── handler.go
│           ├── mocks
│           │   ├── mock_order_assembled_consumer_service.go
│           │   ├── mock_order_paid_consumer_service.go
│           │   └── mock_telegram_service.go
│           ├── service.go
│           └── telegram
│               ├── service.go
│               └── templates
│                   ├── assembled_notification.tmpl
│                   └── paid_notification.tmpl
├── order
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── api
│   │   │   ├── health
│   │   │   │   └── health.go
│   │   │   └── order
│   │   │       └── v1
│   │   │           ├── api.go
│   │   │           ├── cancel.go
│   │   │           ├── create.go
│   │   │           ├── get.go
│   │   │           ├── new_order.go
│   │   │           └── pay.go
│   │   ├── app
│   │   │   ├── app.go
│   │   │   └── di.go
│   │   ├── client
│   │   │   ├── converter
│   │   │   │   └── part.go
│   │   │   └── grpc
│   │   │       ├── client.go
│   │   │       ├── iam
│   │   │       │   └── v1
│   │   │       │       └── client.go
│   │   │       ├── inventory
│   │   │       │   └── v1
│   │   │       │       ├── client.go
│   │   │       │       └── list_parts.go
│   │   │       ├── mocks
│   │   │       │   ├── mock_iam_client.go
│   │   │       │   ├── mock_inventory_client.go
│   │   │       │   └── mock_payment_client.go
│   │   │       └── payment
│   │   │           └── v1
│   │   │               ├── client.go
│   │   │               └── pay_order.go
│   │   ├── config
│   │   │   ├── config.go
│   │   │   ├── env
│   │   │   │   ├── iam_grpc.go
│   │   │   │   ├── inventory_grpc.go
│   │   │   │   ├── kafka.go
│   │   │   │   ├── logger.go
│   │   │   │   ├── order_assembled_consumer.go
│   │   │   │   ├── order_http.go
│   │   │   │   ├── order_paid_producer.go
│   │   │   │   ├── payment_grpc.go
│   │   │   │   └── postgres.go
│   │   │   ├── interfaces.go
│   │   │   └── mocks
│   │   │       ├── mock_inventory_grpc_config.go
│   │   │       ├── mock_kafka_config.go
│   │   │       ├── mock_logger_config.go
│   │   │       ├── mock_order_assembled_consumer_config.go
│   │   │       ├── mock_order_http_config.go
│   │   │       ├── mock_order_paid_producer_config.go
│   │   │       ├── mock_payment_grpc_config.go
│   │   │       └── mock_postgres_config.go
│   │   ├── converter
│   │   │   ├── kafka
│   │   │   │   ├── decoder
│   │   │   │   │   └── order_assembled.go
│   │   │   │   └── kafka.go
│   │   │   ├── order.go
│   │   │   └── user.go
│   │   ├── model
│   │   │   ├── auth.go
│   │   │   ├── error.go
│   │   │   ├── events.go
│   │   │   ├── order.go
│   │   │   ├── part.go
│   │   │   └── user.go
│   │   ├── repository
│   │   │   ├── converter
│   │   │   │   └── order.go
│   │   │   ├── mocks
│   │   │   │   └── mock_order_repository.go
│   │   │   ├── model
│   │   │   │   └── order.go
│   │   │   ├── order
│   │   │   │   ├── create.go
│   │   │   │   ├── get.go
│   │   │   │   ├── repository.go
│   │   │   │   └── update.go
│   │   │   └── repository.go
│   │   └── service
│   │       ├── consumer
│   │       │   └── order_consumer
│   │       │       ├── consumer.go
│   │       │       └── handler.go
│   │       ├── mocks
│   │       │   ├── mock_consumer_service.go
│   │       │   ├── mock_order_producer_service.go
│   │       │   └── mock_order_service.go
│   │       ├── order
│   │       │   ├── cancel.go
│   │       │   ├── cancel_test.go
│   │       │   ├── create.go
│   │       │   ├── create_test.go
│   │       │   ├── get.go
│   │       │   ├── get_test.go
│   │       │   ├── pay.go
│   │       │   ├── pay_test.go
│   │       │   ├── service.go
│   │       │   └── suite_test.go
│   │       ├── producer
│   │       │   └── order_producer
│   │       │       └── producer.go
│   │       └── service.go
│   └── migrations
│       ├── 20250404191615_create_uuid_ossp_extension.sql
│       └── 20250404191624_create_orders_table.sql
├── package-lock.json
├── package.json
├── payment
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   └── internal
│       ├── api
│       │   └── payment
│       │       └── v1
│       │           ├── api.go
│       │           └── pay.go
│       ├── app
│       │   ├── app.go
│       │   └── di.go
│       ├── config
│       │   ├── config.go
│       │   ├── env
│       │   │   ├── logger.go
│       │   │   └── payment_grpc.go
│       │   ├── interfaces.go
│       │   └── mocks
│       │       ├── mock_logger_config.go
│       │       └── mock_payment_grpc_config.go
│       ├── model
│       │   └── errors.go
│       └── service
│           ├── mocks
│           │   └── mock_payment_service.go
│           ├── payment
│           │   ├── pay.go
│           │   ├── pay_test.go
│           │   ├── service.go
│           │   └── suite_test.go
│           └── service.go
├── platform
│   ├── go.mod
│   ├── go.sum
│   └── pkg
│       ├── cache
│       │   ├── client.go
│       │   └── redis
│       │       ├── client.go
│       │       └── set_operator.go
│       ├── closer
│       │   └── closer.go
│       ├── grpc
│       │   └── health
│       │       └── health.go
│       ├── kafka
│       │   ├── consumer
│       │   │   ├── consumer.go
│       │   │   ├── group_handler.go
│       │   │   └── message.go
│       │   ├── kafka.go
│       │   └── producer
│       │       └── producer.go
│       ├── logger
│       │   ├── logger.go
│       │   ├── logger_bench_test.go
│       │   └── noop_logger.go
│       ├── middleware
│       │   ├── grpc
│       │   │   └── auth.go
│       │   ├── http
│       │   │   ├── auth.go
│       │   │   └── error.go
│       │   └── kafka
│       │       └── logging.go
│       ├── migrator
│       │   ├── migrator.go
│       │   └── pg
│       │       └── migrator.go
│       └── testcontainers
│           ├── app
│           │   ├── app.go
│           │   └── opts.go
│           ├── constants.go
│           ├── mongo
│           │   ├── config.go
│           │   ├── connect.go
│           │   ├── init.go
│           │   ├── mongo.go
│           │   └── opts.go
│           ├── network
│           │   └── network.go
│           └── path
│               └── path.go
└── shared
    ├── api
    │   ├── bundles
    │   │   └── order.openapi.v1.bundle.yaml
    │   └── order
    │       └── v1
    │           ├── components
    │           │   ├── create_order_request.yaml
    │           │   ├── create_order_response.yaml
    │           │   ├── enums
    │           │   │   ├── order_status.yaml
    │           │   │   └── payment_method.yaml
    │           │   ├── errors
    │           │   │   ├── bad_gateway_error.yaml
    │           │   │   ├── bad_request_error.yaml
    │           │   │   ├── conflict_error.yaml
    │           │   │   ├── forbidden_error.yaml
    │           │   │   ├── generic_error.yaml
    │           │   │   ├── internal_server_error.yaml
    │           │   │   ├── not_found_error.yaml
    │           │   │   ├── rate_limit_error.yaml
    │           │   │   ├── service_unavailable_error.yaml
    │           │   │   ├── unauthorized_error.yaml
    │           │   │   └── validation_error.yaml
    │           │   ├── get_order_response.yaml
    │           │   ├── order_dto.yaml
    │           │   ├── pay_order_request.yaml
    │           │   └── pay_order_response.yaml
    │           ├── headers
    │           │   └── session_uuid.yaml
    │           ├── order.openapi.yaml
    │           ├── params
    │           │   └── order_uuid.yaml
    │           └── paths
    │               ├── order_by_uuid.yaml
    │               ├── order_cancel.yaml
    │               ├── order_pay.yaml
    │               └── orders.yaml
    ├── go.mod
    ├── go.sum
    ├── pkg
    │   ├── openapi
    │   │   └── order
    │   │       └── v1
    │   │           ├── oas_cfg_gen.go
    │   │           ├── oas_client_gen.go
    │   │           ├── oas_handlers_gen.go
    │   │           ├── oas_interfaces_gen.go
    │   │           ├── oas_json_gen.go
    │   │           ├── oas_labeler_gen.go
    │   │           ├── oas_middleware_gen.go
    │   │           ├── oas_operations_gen.go
    │   │           ├── oas_parameters_gen.go
    │   │           ├── oas_request_decoders_gen.go
    │   │           ├── oas_request_encoders_gen.go
    │   │           ├── oas_response_decoders_gen.go
    │   │           ├── oas_response_encoders_gen.go
    │   │           ├── oas_router_gen.go
    │   │           ├── oas_schemas_gen.go
    │   │           ├── oas_server_gen.go
    │   │           ├── oas_unimplemented_gen.go
    │   │           └── oas_validators_gen.go
    │   └── proto
    │       ├── auth
    │       │   └── v1
    │       │       ├── auth.pb.go
    │       │       └── auth_grpc.pb.go
    │       ├── common
    │       │   └── v1
    │       │       ├── session.pb.go
    │       │       └── user.pb.go
    │       ├── events
    │       │   └── v1
    │       │       └── order.pb.go
    │       ├── inventory
    │       │   └── v1
    │       │       ├── inventory.pb.go
    │       │       └── inventory_grpc.pb.go
    │       ├── payment
    │       │   └── v1
    │       │       ├── payment.pb.go
    │       │       └── payment_grpc.pb.go
    │       └── user
    │           └── v1
    │               ├── user.pb.go
    │               └── user_grpc.pb.go
    └── proto
        ├── auth
        │   └── v1
        │       └── auth.proto
        ├── buf.gen.yaml
        ├── buf.yaml
        ├── common
        │   └── v1
        │       ├── session.proto
        │       └── user.proto
        ├── events
        │   └── v1
        │       └── order.proto
        ├── inventory
        │   └── v1
        │       └── inventory.proto
        ├── payment
        │   └── v1
        │       └── payment.proto
        └── user
            └── v1
                └── user.proto
```

---

### 🔧 Комментарии

- **AuthService — отдельный модуль с `go.mod`**, подключённый в `go.work`.
- **В сервисе выделены слои: `api`, `service`, `repository`, `client`, `config`, `di`** с dependency injection.
- **Middleware для авторизации** работают через gRPC metadata и HTTP заголовки.
- **Сессии передаются между сервисами** через контекст и metadata/заголовки.

---

## ✅ Требования к реализации

### gRPC API:

- Контракт описан в `contracts/iam_service_contracts.md`
- Необходимо реализовать:
    - `Register`
    - `Login`
    - `Whoami`
    - `GetUser`
- Рекомендуется использовать `buf` и сгенерировать код в `shared/pkg/proto`

### PostgreSQL:

- Хранит информацию о пользователях
- Структура таблицы описана в контракте
- Миграции автоматически применяются при старте сервиса

### Redis:

- Используется для хранения сессий
- TTL — 24 часа

### Интеграция с другими сервисами:

- **В `OrderService`** добавить HTTP middleware для проверки авторизации:
  - Использовать заголовок `X-Session-Uuid` для передачи сессии в HTTP запросах
  - Добавить `sessionUUID` в контекст запроса
- **В gRPC клиенте `InventoryService`** внутри `OrderService`:
  - Использовать `ForwardSessionUUIDToGRPC()` для добавления `sessionUUID` в исходящие gRPC metadata
  - Это позволит передавать сессию от Order сервиса к Inventory сервису
- **В `InventoryService`** добавить gRPC interceptor для проверки авторизации:
  - Проверять `session-uuid` в incoming gRPC metadata
  - Валидировать сессию через вызов `IAMService.Whoami`

#### 📖 Incoming vs Outgoing Context

- **Incoming Context** — контекст входящего запроса, содержит metadata/заголовки от клиента
- **Outgoing Context** — контекст исходящего запроса, в который добавляются metadata/заголовки для передачи другим сервисам

---

## 💡 Полезные подсказки

- **Сначала изучите материалы шестой недели**:  
  В уроках шестой недели представлены примеры и объяснения, которые помогут в выполнении задания.  
  👉 Рекомендуется сначала ознакомиться с уроками, а затем приступать к реализации.

- **Не стесняйтесь обращаться за помощью**:  
  Если вы столкнулись с трудностями или у вас возникли вопросы, обращайтесь в чат курса или к своему ревьюеру (если вы на тарифе с проверкой).
  > Если в течение 30 минут вы не можете решить проблему, лучше попросить помощи.  
  > **Спросить — не значит сдаться**, это ускоряет процесс обучения благодаря опыту других.

---

**Автор курса: Олег Козырев, 2025**
