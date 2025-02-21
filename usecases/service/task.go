package service

import (
	"photo_editor/repository"
	"photo_editor/models"
)

type Task struct {
	repo repository.Task
}

func NewTask(repo repository.Task) *Task {
	return &Task{
		repo: repo,
	}
}

func (rs *Task) Get(key string) (*models.Task, error) {
	return rs.repo.Get(key)
}

func (rs *Task) Put(task models.Task) error {
	return rs.repo.Put(task)
}

func (rs *Task) Post(task models.Task) error {
	return rs.repo.Post(task)
}

func (rs *Task) Delete(key string) error {
	return rs.repo.Delete(key)
}