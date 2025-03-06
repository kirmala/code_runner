package repository

import "photo_editor/models"

type User interface {
	Get(key string) (*models.User, error)
	Put(models.User) error
	Post(models.User) error
	Delete(key string) error
}