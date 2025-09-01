package v1

import (
	"context"
	"errors"

	"github.com/space-wanderer/microservices/iam/internal/converter"
	"github.com/space-wanderer/microservices/iam/internal/model"
	userV1 "github.com/space-wanderer/microservices/shared/pkg/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	userUUID, err := a.userService.Register(ctx, converter.ConvertProtoToUserRegistrationInfo(req.GetInfo()))
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists", req.GetInfo())
		}
		return nil, err
	}
	return &userV1.RegisterResponse{UserUuid: userUUID}, nil
}
