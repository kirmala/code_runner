package repository

import (
	"code_processor/http_server/models"
	"context"

	"github.com/google/uuid"
)

type Session interface {
	Get(ctx context.Context, key uuid.UUID) (*models.Session, error)
	Set(ctx context.Context, session models.Session) error
	Delete(ctx context.Context, key uuid.UUID) error
}
