package model

import (
	"time"
)

// Session - модель сессии
type Session struct {
	UUID      string
	UserUUID  string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

// LoginRequest - запрос на вход
type LoginRequest struct {
	Login    string
	Password string
}

// LoginResponse - ответ на вход
type LoginResponse struct {
	SessionUUID string
}

// WhoamiRequest - запрос "кто я"
type WhoamiRequest struct {
	SessionUUID string
}

// WhoamiResponse - ответ "кто я"
type WhoamiResponse struct {
	Session Session
	User    User
}
