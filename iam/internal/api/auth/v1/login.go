package v1

import (
	"context"
	"errors"

	"github.com/space-wanderer/microservices/iam/internal/model"
	authV1 "github.com/space-wanderer/microservices/shared/pkg/proto/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	authResponse, err := a.authService.Login(ctx, model.LoginRequest{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	})
	if err != nil {
		if errors.Is(err, model.ErrUserLoginNotFound) {
			return nil, status.Errorf(codes.NotFound, "user login not found", req.GetLogin())
		}
		return nil, err
	}
	return &authV1.LoginResponse{SessionUuid: authResponse.SessionUUID}, nil
}
