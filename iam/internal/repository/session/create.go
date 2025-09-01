package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/space-wanderer/microservices/iam/internal/repository/model"
)

func (r *repository) CreateSession(ctx context.Context, session model.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	// Сохраняем в Redis с TTL 24 часа
	key := "session:" + session.UUID
	ttl := 24 * time.Hour

	return r.client.SetWithTTL(ctx, key, data, ttl)
}
