package model

import (
	"errors"

	sharedErrors "github.com/space-wanderer/microservices/shared/pkg/errors"
)

var (
	ErrPayment = sharedErrors.NewInvalidArgumentError(errors.New("payment error"))
)
