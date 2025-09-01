package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/space-wanderer/microservices/iam/internal/model"
	repoModel "github.com/space-wanderer/microservices/iam/internal/repository/model"
)

func (s *service) Whoami(ctx context.Context, req model.WhoamiRequest) (model.WhoamiResponse, error) {
	if req.SessionUUID == "" {
		return model.WhoamiResponse{}, fmt.Errorf("invalid session")
	}

	// Получаем сессию из Redis
	session, err := s.sessions.GetSession(ctx, req.SessionUUID)
	if err != nil {
		return model.WhoamiResponse{}, fmt.Errorf("session not found: %w", err)
	}

	// Проверяем, не истекла ли сессия
	if time.Now().UTC().After(session.ExpiresAt) {
		return model.WhoamiResponse{}, fmt.Errorf("session expired")
	}

	// Получаем пользователя по UUID из сессии
	repoUser, err := s.users.GetUserByUUID(ctx, session.UserUUID)
	if err != nil {
		return model.WhoamiResponse{}, fmt.Errorf("user not found: %w", err)
	}

	// Конвертируем в модели сервиса
	modelSession := convertRepoSessionToModelSession(session)
	modelUser := convertRepoUserToModelUser(repoUser)

	return model.WhoamiResponse{
		Session: modelSession,
		User:    modelUser,
	}, nil
}

func convertRepoSessionToModelSession(session repoModel.Session) model.Session {
	return model.Session{
		UUID:      session.UUID,
		UserUUID:  session.UserUUID,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
		ExpiresAt: session.ExpiresAt,
	}
}

func convertRepoUserToModelUser(user repoModel.User) model.User {
	return model.User{
		UUID:      user.UUID,
		Login:     user.Login,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
