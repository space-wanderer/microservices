package user

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/space-wanderer/microservices/iam/internal/repository/model"
)

func (r *repository) CreateUser(ctx context.Context, user model.User) error {
	query := `
		INSERT INTO users (uuid, login, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(ctx, query,
		user.UUID, user.Login, user.Email, user.Password,
		user.CreatedAt, user.UpdatedAt)

	return err
}

func (r *repository) UpdateUser(ctx context.Context, user model.User) error {
	query := `
		UPDATE users
		SET login = $2,
			email = $3,
			password = $4,
			updated_at = $6
		WHERE uuid = $1
	`

	ct, err := r.pool.Exec(ctx, query,
		user.UUID, user.Login, user.Email, user.Password,
		user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *repository) DeleteUser(ctx context.Context, uuid string) error {
	query := `
		DELETE FROM users
		WHERE uuid = $1
	`

	ct, err := r.pool.Exec(ctx, query, uuid)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
