package usecases

import "code_processor/http_server/models"

type User interface {
	Get(key string) (*models.User, error)
	PostLogin(login string, password string) (*string, error)
	PostRegister(models.User) error
	Delete(key string) error
}
