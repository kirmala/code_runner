package basic

import (
	"code_processor/consumer/repository"
	"code_processor/http_server/models"
	"context"
)

type Task struct {
	repo repository.Task
}

func NewTask(repo repository.Task) *Task {
	return &Task{repo: repo}
}

func (t *Task) Put(ctx context.Context, task models.Task) error {
	return t.repo.Put(ctx, task)
}