package repository

import (
	"code_processor/http_server/models"

	"github.com/google/uuid"
)

type Session interface {
	Get(key uuid.UUID) (*models.Session, error)
	Set(models.Session) error
	Delete(key uuid.UUID) error
}
