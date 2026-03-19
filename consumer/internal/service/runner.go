package service

import (
	"context"

	"github.com/kirmala/code_runner/consumer/internal/domain"
)

type Runner interface {
	Run(context.Context, domain.Task) (domain.Task, error)
}
