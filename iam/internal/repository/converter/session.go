package converter

import (
	"github.com/space-wanderer/microservices/iam/internal/model"
	repoModel "github.com/space-wanderer/microservices/iam/internal/repository/model"
)

// ConvertRepoSessionToModelSession конвертирует сессию из репозитория в модель сервиса
func ConvertRepoSessionToModelSession(repoSession repoModel.Session) model.Session {
	return model.Session{
		UUID:      repoSession.UUID,
		UserUUID:  repoSession.UserUUID,
		CreatedAt: repoSession.CreatedAt,
		UpdatedAt: repoSession.UpdatedAt,
		ExpiresAt: repoSession.ExpiresAt,
	}
}

// ConvertModelSessionToRepoSession конвертирует сессию из модели сервиса в репозиторий
func ConvertModelSessionToRepoSession(modelSession model.Session) repoModel.Session {
	return repoModel.Session{
		UUID:      modelSession.UUID,
		UserUUID:  modelSession.UserUUID,
		CreatedAt: modelSession.CreatedAt,
		UpdatedAt: modelSession.UpdatedAt,
		ExpiresAt: modelSession.ExpiresAt,
	}
}
