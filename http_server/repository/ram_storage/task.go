package ram_storage

import (
	"code_processor/http_server/models"
	"code_processor/http_server/repository"

	"github.com/google/uuid"
)

type Task struct {
	data map[uuid.UUID]models.Task
}

func NewTask() *Task {
	return &Task{
		data: make(map[uuid.UUID]models.Task),
	}
}

func (rs Task) Get(key uuid.UUID) (*models.Task, error) {
	value, exists := rs.data[key]
	if !exists {
		return nil, repository.ErrNotFound{Item: "task"}
	}
	return &value, nil
}

func (rs *Task) Put(task models.Task) error {
	_, exists := rs.data[task.Id]
	if !exists {
		return repository.ErrNotFound{Item: "task"}
	}
	rs.data[task.Id] = task
	return nil
}

func (rs *Task) Post(task models.Task) error {
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
