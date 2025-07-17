# План миграции Inventory сервиса на MongoDB

## Обзор
Миграция сервиса inventory с in-memory хранения на MongoDB для обеспечения персистентности данных и масштабируемости.

## Текущее состояние
- **Хранение:** In-memory (sync.Map)
- **Контейнер MongoDB:** Запущен и доступен на localhost:27017
- **База данных:** inventory-service
- **Пользователь:** inventory-service-user / inventory-service-password

## Этап 1: Подготовка MongoDB драйвера

### 1.1 Установка зависимостей
```bash
cd inventory
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/bson
```

### 1.2 Создание конфигурации подключения
- Создать файл `.env` в директории `inventory/`
- Добавить переменные окружения для MongoDB
- Обновить `main.go` для загрузки переменных окружения

## Этап 2: Создание MongoDB репозитория

### 2.1 Структура репозитория
```
inventory/internal/repository/
├── mongo/
│   ├── repository.go      # Интерфейс и основная структура
│   ├── part/
│   │   ├── get.go         # Получение детали по UUID
│   │   ├── list.go        # Список деталей с фильтрацией
│   │   ├── create.go      # Создание новой детали
│   │   ├── update.go      # Обновление детали
│   │   └── delete.go      # Удаление детали
│   └── converter/
│       └── part.go        # Конвертеры между моделями
```

### 2.2 Интерфейс репозитория
```go
type PartRepository interface {
    GetPart(ctx context.Context, uuid string) (*model.Part, error)
    ListParts(ctx context.Context, filter model.PartsFilter) ([]*model.Part, error)
    CreatePart(ctx context.Context, part *model.Part) error
    UpdatePart(ctx context.Context, part *model.Part) error
    DeletePart(ctx context.Context, uuid string) error
}
```

### 2.3 MongoDB структура
```go
type mongoRepository struct {
    client   *mongo.Client
    database *mongo.Database
    parts    *mongo.Collection
}
```

## Этап 3: Реализация MongoDB операций

### 3.1 Подключение к MongoDB
- Создать функцию `NewMongoRepository()`
- Настроить подключение через `mongo.Connect()`
- Добавить health check для проверки соединения

### 3.2 CRUD операции
- **GetPart:** Поиск по UUID с использованием `FindOne()`
- **ListParts:** Поиск с фильтрацией через `Find()` и агрегацию
- **CreatePart:** Вставка документа через `InsertOne()`
- **UpdatePart:** Обновление через `UpdateOne()`
- **DeletePart:** Удаление через `DeleteOne()`

### 3.3 Индексы
```javascript
// Создать индексы для оптимизации запросов
db.parts.createIndex({ "uuid": 1 }, { unique: true })
db.parts.createIndex({ "name": 1 })
db.parts.createIndex({ "category": 1 })
db.parts.createIndex({ "price": 1 })
```

## Этап 4: Конвертеры и модели

### 4.1 BSON теги
```go
type Part struct {
    UUID     string  `bson:"uuid" json:"uuid"`
    Name     string  `bson:"name" json:"name"`
    Category string  `bson:"category" json:"category"`
    Price    float64 `bson:"price" json:"price"`
    // ... другие поля
}
```

### 4.2 Конвертеры
- `ConvertModelToBSON()` - конвертация модели в BSON
- `ConvertBSONToModel()` - конвертация BSON в модель
- `ConvertFilterToBSON()` - конвертация фильтров

## Этап 5: Обновление сервисного слоя

### 5.1 Замена репозитория
- Обновить `service.go` для использования MongoDB репозитория
- Удалить in-memory логику
- Обновить тесты сервиса

### 5.2 Обработка ошибок
- Добавить обработку MongoDB специфичных ошибок
- Реализовать retry логику для сетевых ошибок
- Добавить логирование операций

## Этап 6: Миграция данных

### 6.1 Скрипт миграции
```go
func migrateInMemoryToMongo() {
    // Загрузить данные из in-memory
    // Сохранить в MongoDB
    // Валидировать миграцию
}
```

### 6.2 Инициализация тестовых данных
```go
func seedTestData() {
    // Создать тестовые детали в MongoDB
    // Использовать для тестирования
}
```

## Этап 7: Обновление API слоя

### 7.1 Обновление обработчиков
- Обновить `get.go` и `list.go` в API слое
- Добавить обработку MongoDB ошибок
- Обновить логирование

### 7.2 Валидация
- Добавить валидацию входных данных
- Проверить корректность UUID форматов
- Добавить проверки на существование данных

## Этап 8: Тестирование

### 8.1 Unit тесты
- Создать тесты для MongoDB репозитория
- Использовать `mongo-driver/mongo/integration/mtest`
- Тестировать все CRUD операции

### 8.2 Integration тесты
- Тестировать полный flow от API до MongoDB
- Проверить работу с реальным MongoDB контейнером
- Тестировать обработку ошибок

### 8.3 Performance тесты
- Измерить производительность запросов
- Сравнить с in-memory производительностью
- Оптимизировать индексы при необходимости

## Этап 9: Конфигурация и деплой

### 9.1 Переменные окружения
```env
MONGODB_URI=mongodb://inventory-service-user:inventory-service-password@localhost:27017/inventory-service
MONGODB_DATABASE=inventory-service
MONGODB_COLLECTION=parts
```

### 9.2 Docker конфигурация
- Обновить docker-compose.yml при необходимости
- Добавить health checks для MongoDB
- Настроить volumes для данных

### 9.3 CI/CD
- Обновить GitHub Actions для тестирования с MongoDB
- Добавить тесты с MongoDB в CI pipeline
- Обновить документацию

## Этап 10: Мониторинг и логирование

### 10.1 Метрики
- Добавить метрики для MongoDB операций
- Мониторинг времени выполнения запросов
- Отслеживание ошибок подключения

### 10.2 Логирование
- Логировать все MongoDB операции
- Добавить structured logging
- Настроить уровни логирования

## Этап 11: Документация

### 11.1 API документация
- Обновить OpenAPI спецификацию
- Добавить примеры запросов/ответов
- Документировать новые ошибки

### 11.2 Операционная документация
- Инструкции по развертыванию
- Troubleshooting guide
- Performance tuning guide

## Этап 12: Финальная проверка

### 12.1 Smoke тесты
- Проверить все API endpoints
- Убедиться в корректности данных
- Проверить производительность

### 12.2 Rollback план
- Подготовить план отката к in-memory
- Сохранить резервную копию кода
- Документировать процедуру отката

## Временные рамки
- **Этапы 1-3:** 2-3 дня (базовая структура)
- **Этапы 4-6:** 2-3 дня (реализация и миграция)
- **Этапы 7-9:** 2-3 дня (тестирование и конфигурация)
- **Этапы 10-12:** 1-2 дня (мониторинг и документация)

**Общее время:** 7-11 дней

## Риски и митигация

### Риски
1. **Производительность:** MongoDB может быть медленнее in-memory
2. **Сложность:** Добавление внешней зависимости
3. **Данные:** Риск потери данных при миграции

### Митигация
1. **Производительность:** Оптимизация индексов и запросов
2. **Сложность:** Тщательное тестирование и документация
3. **Данные:** Резервные копии и пошаговая миграция

## Критерии успеха
- [ ] Все API endpoints работают корректно
- [ ] Производительность не ухудшилась более чем на 20%
- [ ] Все тесты проходят
- [ ] Документация обновлена
- [ ] Мониторинг настроен
- [ ] Rollback план готов 