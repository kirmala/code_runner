package basic

import (
	"context"

	"github.com/kirmala/code_runner/consumer/internal/domain"
	"github.com/kirmala/code_runner/consumer/internal/repository"
	"github.com/kirmala/code_runner/consumer/internal/service"
)

type Task struct {
	repo   repository.Task
	runner service.Runner
}

func NewTask(repo repository.Task, runner service.Runner) *Task {
	return &Task{repo: repo, runner: runner}
}

func (ts *Task) Process(ctx context.Context, task domain.Task) error {
	processedTask, err := ts.runner.Run(context.Background(), task)
	if err != nil {
		task.Result = err.Error()
		task.Status = domain.StatusFailed
		_ = ts.repo.Put(context.Background(), task)
		return err
	}
	return ts.repo.Put(context.Background(), processedTask)
}
