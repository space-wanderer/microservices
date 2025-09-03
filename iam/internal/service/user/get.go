package user

import (
	"context"

	"github.com/space-wanderer/microservices/iam/internal/converter"
	"github.com/space-wanderer/microservices/iam/internal/model"
)

func (s *service) GetUser(ctx context.Context, id string) (model.User, error) {
	repoUser, err := s.repo.GetUserByUUID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	return converter.ConvertRepoUserToModelUser(repoUser), nil
}
