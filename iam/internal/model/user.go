package model

import (
	"time"
)

// User - модель пользователя в бизнес-логике
type User struct {
	UUID      string
	Login     string
	Email     string
	Password  string // Хешированный пароль
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserInfo - информация о пользователе для API
type UserInfo struct {
	Login               string
	Email               string
	NotificationMethods []NotificationMethod
}

// NotificationMethod - способ уведомления
type NotificationMethod struct {
	ProviderName string
	Target       string
}

// UserRegistrationInfo - данные для регистрации
type UserRegistrationInfo struct {
	Info     UserInfo
	Password string
}
