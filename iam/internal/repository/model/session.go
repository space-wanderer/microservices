package model

import (
	"time"
)

// Session - модель сессии в репозитории
type Session struct {
	UUID      string
	UserUUID  string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}
