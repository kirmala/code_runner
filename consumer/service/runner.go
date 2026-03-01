package service

import (
	"code_processor/http_server/models"
	"context"
)

type Runner interface {
	Run(context.Context, models.Task) (models.Task, error)
}
