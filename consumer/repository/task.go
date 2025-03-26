package repository

import "code_processor/consumer/models"

type Task interface {
	Put(models.Task) error
}
