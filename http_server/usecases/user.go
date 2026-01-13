package usecases

import (
	"code_processor/http_server/models"

	"github.com/google/uuid"
)

type User interface {
	Get(key uuid.UUID) (*models.User, error)
	PostLogin(login string, password string) (uuid.UUID, error)
	PostRegister(models.User) error
	Delete(key uuid.UUID) error
}
