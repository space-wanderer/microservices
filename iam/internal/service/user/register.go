package user

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/space-wanderer/microservices/iam/internal/converter"
	"github.com/space-wanderer/microservices/iam/internal/model"
)

func (s *service) Register(ctx context.Context, info model.UserRegistrationInfo) (string, error) {
	if info.Password == "" || info.Info.Login == "" || info.Info.Email == "" {
		return "", fmt.Errorf("invalid registration data")
	}

	// Проверка уникальности логина/почты (минимально)
	if _, err := s.repo.GetUserByLogin(ctx, info.Info.Login); err == nil {
		return "", fmt.Errorf("login already exists")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}
	if _, err := s.repo.GetUserByEmail(ctx, info.Info.Email); err == nil {
		return "", fmt.Errorf("email already exists")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}

	uuid, err := generateUUID()
	if err != nil {
		return "", err
	}

	user := converter.ConvertUserInfoToModelUser(info.Info)
	user.Password = hashPassword(info.Password)
	user.UUID = uuid
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt

	repoUser := converter.ConvertModelUserToRepoUser(user)
	if err := s.repo.CreateUser(ctx, repoUser); err != nil {
		return "", err
	}

	return user.UUID, nil
}

func hashPassword(password string) string {
	// Пример: sha256. Для продакшена использовать bcrypt/argon2.
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
