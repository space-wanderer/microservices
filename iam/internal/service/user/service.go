package user

import (
	"github.com/space-wanderer/microservices/iam/internal/repository"
)

type service struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) *service {
	return &service{repo: repo}
}
