package usecases

import "code_processor/http_server/models"

type CodeProcessor interface {
	Process(models.Task) (*models.Task, error)
}
