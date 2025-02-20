package usecases

import "photo_editor/models"

type Task interface {
	Get(key string) (*models.Task, error)
	Put(models.Task) error
	Post(models.Task) error
	Delete(key string) error
}