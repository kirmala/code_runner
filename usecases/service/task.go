package service

import (
	"code_runner/repository"
	"code_runner/models"
	"time"
)

type Task struct {
	taskRepo repository.Task
	sessionRepo repository.Session
}

func edit(service Task, newTask *models.Task) {
	time.Sleep(8 * time.Second)
	newTask.Status = "ready"
	newTask.Result = "something happend"
	service.Put(*newTask)
}

func NewTask(taskRepo repository.Task, sessionRepo repository.Session) *Task {
	return &Task{
		taskRepo: taskRepo,
		sessionRepo: sessionRepo,
	}
}

func (rs *Task) GetUserId(key string) (*string, error) {
	if key == "" {
		return nil, repository.NotFound
	}
	session, err := rs.sessionRepo.Get(key)
	if (err != nil) {
		return nil, err
	}
	return &session.UserId, err
}

func (rs *Task) GetStatus(key string) (*string, error) {
	task, err := rs.taskRepo.Get(key)
	if (err != nil) {
		return nil, err
	}
	return &task.Status, err
}

func (rs *Task) GetResult(key string) (*string, error) {
	task, err := rs.taskRepo.Get(key)
	if (err != nil) {
		return nil, err
	}
	return &task.Result, err
}


func (rs *Task) Put(task models.Task) error {
	return rs.taskRepo.Put(task)
}

func (rs *Task) Post(task models.Task) error {
	go edit(*rs, &task)
	return rs.taskRepo.Post(task)
}

func (rs *Task) Delete(key string) error {
	return rs.taskRepo.Delete(key)
}