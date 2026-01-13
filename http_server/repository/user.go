package repository

import (
	"code_processor/http_server/models"

	"github.com/google/uuid"
)

type User interface {
	GetByLogin(login string) (*models.User, error)
	GetById(id uuid.UUID) (*models.User, error)
	Put(models.User) error
	Post(models.User) error
	Delete(key uuid.UUID) error
}
