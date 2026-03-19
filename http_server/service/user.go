package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/domain"
)

type User interface {
	Login(ctx context.Context, login string, password string) (uuid.UUID, error)
	Register(domain.User) error
}
