package auth

import (
	"github.com/space-wanderer/microservices/iam/internal/repository"
)

type service struct {
	users    repository.UserRepository
	sessions repository.SessionRepository
}

func NewService(users repository.UserRepository, sessions repository.SessionRepository) *service {
	return &service{
		users:    users,
		sessions: sessions,
	}
}
