package repository

import "code_processor/http_server/models"

type TaskSender interface {
	Send(models.Task) error
}
