package model

import (
	"time"
)

// User - модель пользователя в репозитории
type User struct {
	UUID      string
	Login     string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
