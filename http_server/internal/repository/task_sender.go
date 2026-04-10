package repository

import (
	"context"

	"github.com/kirmala/code_runner/http_server/internal/domain"
)

type TaskSender interface {
	Send(context.Context, domain.Task) error
}
