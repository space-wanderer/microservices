package model

import (
	"errors"

	sharedErrors "github.com/space-wanderer/microservices/shared/pkg/errors"
)

var (
	ErrUserNotFound      = sharedErrors.NewNotFoundError(errors.New("user not found"))
	ErrUserAlreadyExists = sharedErrors.NewInvalidArgumentError(errors.New("user already exists"))
	ErrUserLoginNotFound = sharedErrors.NewNotFoundError(errors.New("user login not found"))
	ErrSessionNotFound   = sharedErrors.NewNotFoundError(errors.New("session not found"))
)
