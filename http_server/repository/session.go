package repository

import "code_processor/http_server/models"

type Session interface {
	Get(key string) (*models.Session, error)
	Set(models.Session) error
	Delete(key string) error
}
