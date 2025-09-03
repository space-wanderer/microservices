package session

import (
	"context"

	"github.com/space-wanderer/microservices/iam/internal/repository/model"
)

func (r *repository) AddSessionToUserSet(ctx context.Context, session model.Session) error {
	return r.client.SAdd(ctx, "session:"+session.UUID, session.UserUUID)
}
