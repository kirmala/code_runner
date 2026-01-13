package repository

import (
	"code_processor/http_server/models"

	"github.com/google/uuid"
)

type Task interface {
	Get(key uuid.UUID) (*models.Task, error)
	Put(models.Task) error
	Post(models.Task) error
	Delete(key uuid.UUID) error
}
