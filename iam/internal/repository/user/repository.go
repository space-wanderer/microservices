package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	repoiface "github.com/space-wanderer/microservices/iam/internal/repository"
)

var _ repoiface.UserRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}
