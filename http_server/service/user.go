package service

import (
	"code_processor/http_server/models"
	"context"

	"github.com/google/uuid"
)

type User interface {
	Login(ctx context.Context, login string, password string) (uuid.UUID, error)
	Register(models.User) error
}
