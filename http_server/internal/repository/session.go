package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/internal/domain"
)

type Session interface {
	Get(ctx context.Context, key uuid.UUID) (*domain.Session, error)
	Set(ctx context.Context, session domain.Session) error
	Delete(ctx context.Context, key uuid.UUID) error
}
