package v1

import (
	"context"
	"errors"

	"github.com/space-wanderer/microservices/iam/internal/converter"
	"github.com/space-wanderer/microservices/iam/internal/model"
	authV1 "github.com/space-wanderer/microservices/shared/pkg/proto/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	authResponse, err := a.authService.Whoami(ctx, model.WhoamiRequest{SessionUUID: req.GetSessionUuid()})
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			return nil, status.Errorf(codes.NotFound, "session not found", req.GetSessionUuid())
		}
		return nil, err
	}
	return &authV1.WhoamiResponse{
		Session: converter.ConvertModelSessionToProto(authResponse.Session),
		User:    converter.ConvertModelUserToProto(authResponse.User),
	}, nil
}
