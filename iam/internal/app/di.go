package app

import (
	"context"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"

	authV1API "github.com/space-wanderer/microservices/iam/internal/api/auth/v1"
	userV1API "github.com/space-wanderer/microservices/iam/internal/api/user/v1"
	"github.com/space-wanderer/microservices/iam/internal/config"
	"github.com/space-wanderer/microservices/iam/internal/repository"
	sessionRepository "github.com/space-wanderer/microservices/iam/internal/repository/session"
	userRepository "github.com/space-wanderer/microservices/iam/internal/repository/user"
	"github.com/space-wanderer/microservices/iam/internal/service"
	authService "github.com/space-wanderer/microservices/iam/internal/service/auth"
	userService "github.com/space-wanderer/microservices/iam/internal/service/user"
	"github.com/space-wanderer/microservices/platform/pkg/cache"
	"github.com/space-wanderer/microservices/platform/pkg/cache/redis"
	"github.com/space-wanderer/microservices/platform/pkg/closer"
	"github.com/space-wanderer/microservices/platform/pkg/logger"
	authV1 "github.com/space-wanderer/microservices/shared/pkg/proto/auth/v1"
	userV1 "github.com/space-wanderer/microservices/shared/pkg/proto/user/v1"
)

type diContainer struct {
	userV1API userV1.UserServiceServer
	authV1API authV1.AuthServiceServer

	userService service.UserService
	authService service.AuthService

	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository

	pgPool      *pgxpool.Pool
	redisClient cache.RedisClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userV1API.NewAPI(d.UserService(ctx))
	}
	return d.userV1API
}

func (d *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = authV1API.NewAPI(d.AuthService(ctx))
	}
	return d.authV1API
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewService(d.UserRepository(ctx))
	}
	return d.userService
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewService(
			d.UserRepository(ctx),
			d.SessionRepository(),
		)
	}
	return d.authService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewRepository(d.PGPool(ctx))
	}
	return d.userRepository
}

func (d *diContainer) SessionRepository() repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepository.NewRepository(d.RedisClient())
	}
	return d.sessionRepository
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to create PostgreSQL pool: %s", err.Error()))
		}

		// Проверяем соединение
		if err := pool.Ping(ctx); err != nil {
			panic(fmt.Sprintf("failed to ping PostgreSQL: %v", err))
		}

		closer.AddNamed("PostgreSQL pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}
	return d.pgPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		redisConfig := config.AppConfig().Redis
		// Создаем Redis pool
		pool := &redigo.Pool{
			MaxIdle:     redisConfig.MaxIdle(),
			IdleTimeout: redisConfig.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", redisConfig.Address())
			},
		}

		d.redisClient = redis.NewClient(pool, &logger.NoopLogger{}, redisConfig.ConnectionTimeout())
	}
	return d.redisClient
}
