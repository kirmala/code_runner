package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/internal/domain"
)

type Task interface {
	GetStatus(key uuid.UUID) (string, error)
	GetResult(key uuid.UUID) (*string, error)
	Put(domain.Task) error
	Post(context.Context, domain.Task) error
	Delete(key uuid.UUID) error
}
