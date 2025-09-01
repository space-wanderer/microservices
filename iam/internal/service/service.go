package service

import (
	"context"

	"github.com/space-wanderer/microservices/iam/internal/model"
)

type UserService interface {
	Register(ctx context.Context, info model.UserRegistrationInfo) (string, error)
	GetUser(ctx context.Context, uuid string) (model.User, error)
}

type AuthService interface {
	Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error)
	Whoami(ctx context.Context, req model.WhoamiRequest) (model.WhoamiResponse, error)
}
