package usecases

import (
	"code_processor/http_server/models"

	"github.com/google/uuid"
)

type Task interface {
	GetStatus(key uuid.UUID) (string, error)
	GetResult(key uuid.UUID) (*string, error)
	GetUserId(key uuid.UUID) (uuid.UUID, error)
	Put(models.Task) error
	Post(models.Task) error
	Delete(key uuid.UUID) error
}
