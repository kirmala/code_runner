package usecases

import "consumer/models"

type CodeProcessor interface {
	Process(models.Task) (*models.Task, error)
}
