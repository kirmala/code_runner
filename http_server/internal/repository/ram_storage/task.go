package ram_storage

import (
	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/internal/domain"
	"github.com/kirmala/code_runner/http_server/internal/repository"
)

type Task struct {
	data map[uuid.UUID]domain.Task
}

func NewTask() *Task {
	return &Task{
		data: make(map[uuid.UUID]domain.Task),
	}
}

func (rs Task) Get(key uuid.UUID) (*domain.Task, error) {
	value, exists := rs.data[key]
	if !exists {
		return nil, repository.ErrNotFound{Item: "task"}
	}
	return &value, nil
}

func (rs *Task) Put(task domain.Task) error {
	_, exists := rs.data[task.Id]
	if !exists {
		return repository.ErrNotFound{Item: "task"}
	}
	rs.data[task.Id] = task
	return nil
}

func (rs *Task) Post(task domain.Task) error {
	if _, exists := rs.data[task.Id]; exists {
		return repository.ErrConflict{Field: "id"}
	}
	rs.data[task.Id] = task
	return nil
}

func (rs *Task) Delete(key uuid.UUID) error {
	if _, exists := rs.data[key]; !exists {
		return repository.ErrNotFound{Item: "task"}
	}
	delete(rs.data, key)
	return nil
}
