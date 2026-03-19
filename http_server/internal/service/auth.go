package service

import (
	"context"

	"github.com/google/uuid"
)

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (uuid.UUID, error)
}