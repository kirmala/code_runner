package repository

import "code_runner/models"

type Task interface {
	Get(key string) (*models.Task, error)
	Put(models.Task) error
	Post(models.Task) error
	Delete(key string) error
}