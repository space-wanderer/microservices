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

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	user, err := a.userService.GetUser(ctx, req.GetUserUuid())
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found", req.GetUserUuid())
		}
		return nil, err
	}
	return &userV1.GetUserResponse{
		User: converter.ConvertModelUserToProto(user),
	}, nil
}
