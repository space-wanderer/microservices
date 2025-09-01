package session

import "github.com/space-wanderer/microservices/platform/pkg/cache"

type repository struct {
	client cache.RedisClient
}

func NewRepository(client cache.RedisClient) *repository {
	return &repository{client: client}
}
