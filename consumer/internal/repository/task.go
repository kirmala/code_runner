package repository

import (
	"context"

	"github.com/kirmala/code_runner/consumer/internal/domain"
)

type Task interface {
	Put(context.Context, domain.Task) error
}
