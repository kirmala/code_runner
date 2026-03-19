package service

import (
	"context"

	"github.com/kirmala/code_runner/consumer/internal/domain"
)

type Task interface {
	Put(ctx context.Context, task domain.Task) error
}