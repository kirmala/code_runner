package repository

import "code_processor/http_server/models"

type Task interface {
	Put(models.Task) error
}
