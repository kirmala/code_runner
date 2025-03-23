package repository

import "code_runner/models"

type TaskSender interface {
	Send(models.Task) error
}