package service

import (
	"photo_editor/repository"
	"photo_editor/models"
	"time"
)

type Task struct {
	repo repository.Task
}

func edit(service Task, newTask *models.Task) {
	time.Sleep(8 * time.Second)
	newTask.Status = "ready"
	newTask.Result = "something happend"
	service.Put(*newTask)
}

func NewTask(repo repository.Task) *Task {
	return &Task{
		repo: repo,
	}
}

func (rs *Task) GetStatus(key string) (*string, error) {
	task, err := rs.repo.Get(key)
	if (err != nil) {
		return nil, err
	}
	return &task.Status, err
}

func (rs *Task) GetResult(key string) (*string, error) {
	task, err := rs.repo.Get(key)
	if (err != nil) {
		return nil, err
	}
	return &task.Result, err
}


func (rs *Task) Put(task models.Task) error {
	return rs.repo.Put(task)
}

func (rs *Task) Post(task models.Task) error {
	go edit(*rs, &task)
	return rs.repo.Post(task)
}

func (rs *Task) Delete(key string) error {
	return rs.repo.Delete(key)
}