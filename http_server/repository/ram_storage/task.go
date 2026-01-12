package ram_storage

import (
	"code_processor/http_server/models"
	"code_processor/http_server/repository"
)

type Task struct {
	data map[string]models.Task
}

func NewTask() *Task {
	return &Task{
		data: make(map[string]models.Task),
	}
}

func (rs *Task) Get(key string) (*models.Task, error) {
	value, exists := rs.data[key]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return &value, nil
}

func (rs *Task) Put(task models.Task) error {
	rs.data[task.Id] = task
	return nil
}

func (rs *Task) Post(task models.Task) error {
	if _, exists := rs.data[task.Id]; exists {
		return repository.ErrAlreadyExists
	}
	rs.data[task.Id] = task
	return nil
}

func (rs *Task) Delete(key string) error {
	if _, exists := rs.data[key]; !exists {
		return repository.ErrNotFound
	}
	delete(rs.data, key)
	return nil
}
