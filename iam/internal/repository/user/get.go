package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	repoerrors "github.com/space-wanderer/microservices/iam/internal/model"
	"github.com/space-wanderer/microservices/iam/internal/repository/model"
)

func (r *repository) GetUserByUUID(ctx context.Context, uuid string) (model.User, error) {
	query := `
		SELECT uuid, login, email, password, created_at, updated_at
		FROM users
		WHERE uuid = $1
	`

	var user model.User
	err := r.pool.QueryRow(ctx, query, uuid).Scan(
		&user.UUID, &user.Login, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, repoerrors.ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (r *repository) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	query := `
		SELECT uuid, login, email, password, created_at, updated_at
		FROM users
		WHERE login = $1
	`

	var user model.User
	err := r.pool.QueryRow(ctx, query, login).Scan(
		&user.UUID, &user.Login, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, repoerrors.ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	query := `
		SELECT uuid, login, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user model.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.UUID, &user.Login, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, repoerrors.ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}
