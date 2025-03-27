package repository

import "code_processor/http_server/models"

type Task interface {
	Get(key string) (*models.Task, error)
	Put(models.Task) error
	Post(models.Task) error
	Delete(key string) error
}
