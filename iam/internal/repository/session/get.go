package session

import (
	"context"
	"encoding/json"

	"github.com/space-wanderer/microservices/iam/internal/repository/model"
)

func (r *repository) GetSession(ctx context.Context, uuid string) (model.Session, error) {
	key := "session:" + uuid

	data, err := r.client.Get(ctx, key)
	if err != nil {
		return model.Session{}, err
	}

	var session model.Session
	err = json.Unmarshal(data, &session)
	if err != nil {
		return model.Session{}, err
	}

	return session, nil
}
