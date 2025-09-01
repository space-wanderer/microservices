package converter

import (
	"github.com/space-wanderer/microservices/iam/internal/model"
	authV1 "github.com/space-wanderer/microservices/shared/pkg/proto/auth/v1"
	commonV1 "github.com/space-wanderer/microservices/shared/pkg/proto/common/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertModelLoginRequestToProto конвертирует запрос входа из модели в proto
func ConvertModelLoginRequestToProto(req model.LoginRequest) *authV1.LoginRequest {
	return &authV1.LoginRequest{
		Login:    req.Login,
		Password: req.Password,
	}
}

// ConvertProtoLoginResponseToModel конвертирует ответ входа из proto в модель
func ConvertProtoLoginResponseToModel(resp *authV1.LoginResponse) model.LoginResponse {
	return model.LoginResponse{
		SessionUUID: resp.SessionUuid,
	}
}

// ConvertModelWhoamiRequestToProto конвертирует запрос whoami из модели в proto
func ConvertModelWhoamiRequestToProto(req model.WhoamiRequest) *authV1.WhoamiRequest {
	return &authV1.WhoamiRequest{
		SessionUuid: req.SessionUUID,
	}
}

// ConvertProtoWhoamiResponseToModel конвертирует ответ whoami из proto в модель
func ConvertProtoWhoamiResponseToModel(resp *authV1.WhoamiResponse) model.WhoamiResponse {
	return model.WhoamiResponse{
		Session: ConvertProtoSessionToModel(resp.Session),
		User:    ConvertProtoUserToModel(resp.User),
	}
}

// ConvertModelSessionToProto конвертирует сессию из модели в proto
func ConvertModelSessionToProto(session model.Session) *commonV1.Session {
	return &commonV1.Session{
		Uuid:      session.UUID,
		CreatedAt: timestamppb.New(session.CreatedAt),
		UpdatedAt: timestamppb.New(session.UpdatedAt),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}

// ConvertProtoSessionToModel конвертирует сессию из proto в модель
func ConvertProtoSessionToModel(session *commonV1.Session) model.Session {
	return model.Session{
		UUID:      session.Uuid,
		CreatedAt: session.CreatedAt.AsTime(),
		UpdatedAt: session.UpdatedAt.AsTime(),
		ExpiresAt: session.ExpiresAt.AsTime(),
	}
}

// ConvertModelUserToProto конвертирует пользователя из модели в proto
func ConvertModelUserToProto(user model.User) *commonV1.User {
	return &commonV1.User{
		Uuid:      user.UUID,
		Info:      ConvertModelUserInfoToProto(user),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

// ConvertProtoUserToModel конвертирует пользователя из proto в модель
func ConvertProtoUserToModel(user *commonV1.User) model.User {
	return model.User{
		UUID:      user.Uuid,
		Login:     user.Info.Login,
		Email:     user.Info.Email,
		CreatedAt: user.CreatedAt.AsTime(),
		UpdatedAt: user.UpdatedAt.AsTime(),
	}
}

// ConvertModelUserInfoToProto конвертирует информацию пользователя из модели в proto
func ConvertModelUserInfoToProto(user model.User) *commonV1.UserInfo {
	// Пока возвращаем пустой массив notification methods
	// В будущем можно добавить отдельную таблицу для хранения методов уведомлений
	return &commonV1.UserInfo{
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: []*commonV1.NotificationMethod{},
	}
}
