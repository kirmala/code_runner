package repository

import "code_runner/models"

type Session interface {
	Get(key string) (*models.Session, error)
	Post(models.Session) error
	Delete(key string) error
}