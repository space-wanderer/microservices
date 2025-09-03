package converter

import (
	"github.com/space-wanderer/microservices/iam/internal/model"
	repoModel "github.com/space-wanderer/microservices/iam/internal/repository/model"
)

// ConvertRepoUserToModelUser конвертирует пользователя из репозитория в модель сервиса
func ConvertRepoUserToModelUser(repoUser repoModel.User) model.User {
	return model.User{
		UUID:      repoUser.UUID,
		Login:     repoUser.Login,
		Email:     repoUser.Email,
		Password:  repoUser.Password,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: repoUser.UpdatedAt,
	}
}

// ConvertModelUserToRepoUser конвертирует пользователя из модели сервиса в репозиторий
func ConvertModelUserToRepoUser(modelUser model.User) repoModel.User {
	return repoModel.User{
		UUID:      modelUser.UUID,
		Login:     modelUser.Login,
		Email:     modelUser.Email,
		Password:  modelUser.Password,
		CreatedAt: modelUser.CreatedAt,
		UpdatedAt: modelUser.UpdatedAt,
	}
}
