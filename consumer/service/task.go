package service

import (
	"code_processor/http_server/models"
	"context"
)

type Task interface {
	Put(ctx context.Context, task models.Task) error
}