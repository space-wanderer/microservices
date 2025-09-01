package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/space-wanderer/microservices/iam/internal/model"
	repoModel "github.com/space-wanderer/microservices/iam/internal/repository/model"
)

func (s *service) Login(ctx context.Context, req model.LoginRequest) (model.LoginResponse, error) {
	if req.Login == "" || req.Password == "" {
		return model.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	u, err := s.users.GetUserByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.LoginResponse{}, fmt.Errorf("invalid credentials")
		}
		return model.LoginResponse{}, err
	}

	if u.Password != hashPassword(req.Password) {
		return model.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	// Создаем сессию
	sessionUUID, err := generateUUID()
	if err != nil {
		return model.LoginResponse{}, err
	}

	now := time.Now().UTC()
	session := model.Session{
		UUID:      sessionUUID,
		UserUUID:  u.UUID,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(24 * time.Hour),
	}

	// Конвертируем в модель репозитория и сохраняем в Redis
	repoSession := convertModelSessionToRepoSession(session)
	if err := s.sessions.CreateSession(ctx, repoSession); err != nil {
		return model.LoginResponse{}, fmt.Errorf("failed to create session: %w", err)
	}

	return model.LoginResponse{SessionUUID: sessionUUID}, nil
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func generateUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// RFC 4122 variant 10 and version 4
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

func convertModelSessionToRepoSession(session model.Session) repoModel.Session {
	return repoModel.Session{
		UUID:      session.UUID,
		UserUUID:  session.UserUUID,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
		ExpiresAt: session.ExpiresAt,
	}
}
