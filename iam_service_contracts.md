## 🔐 gRPC контракты IAM: `auth.v1` и `user.v1`

Набор сервисов для аутентификации и управления пользователями. Два gRPC сервиса: `AuthService` (пакет `auth.v1`) и `UserService` (пакет `user.v1`). Общие структуры вынесены в пакет `common.v1`.

---

### Сервис: AuthService (`auth.v1`)

Методы:

1) Login(LoginRequest) → LoginResponse

LoginRequest

| Поле   | Тип    | Описание            |
|--------|--------|---------------------|
| login  | string | Логин пользователя  |
| password | string | Пароль пользователя |

LoginResponse

| Поле         | Тип    | Описание                 |
|--------------|--------|--------------------------|
| session_uuid | string | UUID активной сессии     |

2) Whoami(WhoamiRequest) → WhoamiResponse

WhoamiRequest

| Поле         | Тип    | Описание                 |
|--------------|--------|--------------------------|
| session_uuid | string | UUID активной сессии     |

WhoamiResponse

| Поле    | Тип                | Описание                          |
|---------|--------------------|-----------------------------------|
| session | common.v1.Session  | Информация о текущей сессии       |
| user    | common.v1.User     | Владелец текущей сессии           |

---

### Сервис: UserService (`user.v1`)

Методы:

1) Register(RegisterRequest) → RegisterResponse

UserRegistrationInfo

| Поле  | Тип                | Описание                       |
|-------|--------------------|--------------------------------|
| info  | common.v1.UserInfo | Основная информация пользователя |
| password | string          | Пароль                         |

RegisterRequest

| Поле | Тип                    | Описание                 |
|------|------------------------|--------------------------|
| info | UserRegistrationInfo   | Данные для регистрации   |

RegisterResponse

| Поле      | Тип    | Описание                         |
|-----------|--------|----------------------------------|
| user_uuid | string | UUID созданного пользователя     |

2) GetUser(GetUserRequest) → GetUserResponse

GetUserRequest

| Поле      | Тип    | Описание               |
|-----------|--------|------------------------|
| user_uuid | string | UUID пользователя      |

GetUserResponse

| Поле | Тип            | Описание        |
|------|----------------|-----------------|
| user | common.v1.User | Пользователь    |

---

### Общие структуры (`common.v1`)

Session

| Поле       | Тип                         | Описание                               |
|------------|-----------------------------|----------------------------------------|
| uuid       | string                      | UUID сессии                            |
| created_at | google.protobuf.Timestamp   | Время создания                         |
| updated_at | google.protobuf.Timestamp   | Время последнего обновления            |
| expires_at | google.protobuf.Timestamp   | Время истечения                        |

NotificationMethod

| Поле          | Тип    | Описание                                           |
|---------------|--------|----------------------------------------------------|
| provider_name | string | Провайдер: `telegram`, `email`, `push` и т.д.     |
| target        | string | Адрес/идентификатор назначения (email, чат-id)     |

UserInfo

| Поле                 | Тип                      | Описание                     |
|----------------------|--------------------------|------------------------------|
| login                | string                   | Логин                        |
| email                | string                   | Email                        |
| notification_methods | []NotificationMethod     | Каналы уведомлений           |

User

| Поле       | Тип              | Описание                   |
|------------|------------------|----------------------------|
| uuid       | string           | UUID пользователя          |
| info       | UserInfo         | Базовая информация         |
| created_at | google.protobuf.Timestamp | Дата создания    |
| updated_at | google.protobuf.Timestamp | Дата обновления  |

