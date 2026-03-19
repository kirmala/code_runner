package basic

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/domain"
	"github.com/kirmala/code_runner/http_server/repository"
	//"time"
)

type Task struct {
	taskRepo    repository.Task
	taskSender  repository.TaskSender
}

// func edit(service Task, newTask *domain.Task) {
// 	time.Sleep(8 * time.Second)
// 	newTask.Status = "ready"
// 	newTask.Result = "something happend"
// 	service.Put(*newTask)
// }

func NewTask(taskRepo repository.Task, sessionRepo repository.Session, taskSender repository.TaskSender) *Task {
	return &Task{
		taskRepo:    taskRepo,
		taskSender:  taskSender,
	}
}

func (rs *Task) GetStatus(key uuid.UUID) (string, error) {
	task, err := rs.taskRepo.Get(key)
	if err != nil {
		return "", err
	}
	return task.Status.String(), err
}

func (rs *Task) GetResult(key uuid.UUID) (*string, error) {
	task, err := rs.taskRepo.Get(key)
	if err != nil {
		return nil, err
	}
	return &task.Result, err
}

func (rs *Task) Put(task domain.Task) error {
	return rs.taskRepo.Put(task)
}

func (rs *Task) Post(task domain.Task) error {
	err := rs.taskSender.Send(task)
	if err != nil {
		return fmt.Errorf("sending task: %w", err)
	}
	return rs.taskRepo.Post(task)
}

func (rs *Task) Delete(key uuid.UUID) error {
	return rs.taskRepo.Delete(key)
}
