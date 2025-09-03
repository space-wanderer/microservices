package converter

import (
	"github.com/space-wanderer/microservices/iam/internal/model"
	repoModel "github.com/space-wanderer/microservices/iam/internal/repository/model"
	userV1 "github.com/space-wanderer/microservices/shared/pkg/proto/user/v1"
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

// ConvertModelUserToUserInfo конвертирует пользователя в информацию для API
func ConvertModelUserToUserInfo(modelUser model.User) model.UserInfo {
	return model.UserInfo{
		Login:               modelUser.Login,
		Email:               modelUser.Email,
		NotificationMethods: []model.NotificationMethod{}, // Пока пустой массив
	}
}

// ConvertUserInfoToModelUser конвертирует информацию пользователя в модель
func ConvertUserInfoToModelUser(userInfo model.UserInfo) model.User {
	return model.User{
		Login: userInfo.Login,
		Email: userInfo.Email,
		// NotificationMethods не сохраняются в основной модели User
	}
}

// ConvertProtoToUserRegistrationInfo конвертирует proto в UserRegistrationInfo
func ConvertProtoToUserRegistrationInfo(info *userV1.UserRegistrationInfo) model.UserRegistrationInfo {
	return model.UserRegistrationInfo{
		Info: model.UserInfo{
			Login:               info.GetInfo().GetLogin(),
			Email:               info.GetInfo().GetEmail(),
			NotificationMethods: []model.NotificationMethod{}, // Пока пустой
		},
		Password: info.GetPassword(),
	}
}
