package service

import (
	"code_processor/http_server/models"

	"github.com/google/uuid"
)

type User interface {
	Login(login string, password string) (uuid.UUID, error)
	Register(models.User) error
}
