package repository

import (
	"code_processor/http_server/models"
	"context"
)

type Task interface {
	Put(context.Context, models.Task) error
}
