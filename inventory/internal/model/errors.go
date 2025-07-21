package model

import (
	"errors"

	sharedErrors "github.com/space-wanderer/microservices/shared/pkg/errors"
)

var (
	ErrPartNotFound = sharedErrors.NewNotFoundError(errors.New("part not found"))
	ErrInvalidUUID  = sharedErrors.NewInvalidArgumentError(errors.New("invalid uuid"))
)
