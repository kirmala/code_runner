package repository

import "photo_editor/models"

type Session interface {
	Get(key string) (*models.Session, error)
	Post(models.Session) error
	Delete(key string) error
}