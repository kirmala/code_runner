package usecases

import "code_processor/consumer/models"

type CodeProcessor interface {
	Process(models.Task) (*models.Task, error)
}
