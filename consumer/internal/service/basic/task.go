package basic

import (
	"context"

	"github.com/kirmala/code_runner/consumer/internal/domain"
	"github.com/kirmala/code_runner/consumer/internal/repository"
)

type Task struct {
	repo repository.Task
}

func NewTask(repo repository.Task) *Task {
	return &Task{repo: repo}
}

func (t *Task) Put(ctx context.Context, task domain.Task) error {
	return t.repo.Put(ctx, task)
}