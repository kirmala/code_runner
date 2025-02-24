package usecases

import "photo_editor/models"

type Task interface {
	GetStatus(key string) (*string, error)
	GetResult(key string) (*string, error)
	Put(models.Task) error
	Post(models.Task) error
	Delete(key string) error
}