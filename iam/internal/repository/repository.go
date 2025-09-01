package repository

import (
	"context"

	"github.com/space-wanderer/microservices/iam/internal/repository/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUserByUUID(ctx context.Context, uuid string) (model.User, error)
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, uuid string) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session model.Session) error
	GetSession(ctx context.Context, uuid string) (model.Session, error)
}
