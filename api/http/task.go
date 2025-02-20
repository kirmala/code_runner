package http

import (
	"net/http"
	"photo_editor/api/http/types"
	"photo_editor/editing"
	"photo_editor/models"
	"photo_editor/usecases"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Task struct {
	service usecases.Task
}

func NewHandler(service usecases.Task) *Task {
	return &Task{service: service}
}

// @Summary Get a task status
// @Description Get a task status by its id
// @Tags task
// @Accept  json
// @Produce json
// @Param id path string true "Id of the task"
// @Success 200 {task} types.GetTaskStatusHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Task not found"
// @Router /task/status/{id} [get]
func (s *Task) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	task, err := s.service.Get(req.Id)
	types.ProcessError(w, err, &types.GetTaskStatusHandlerResponse{Status: &task.Status})
}

// @Summary Get a task result
// @Description Get a task result by its id
// @Tags task
// @Accept  json
// @Produce json
// @Param id path string true "Id of the task"
// @Success 200 {task} types.GetTaskResultHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Task not found"
// @Router /task/result/{id} [get]
func (s *Task) getResultHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	task, err := s.service.Get(req.Id)
	types.ProcessError(w, err, &types.GetTaskResultHandlerResponse{Result: &task.Result})
}


// @Summary Create an task
// @Description Create a new task with the specified key and value
// @Tags task
// @Accept  json
// @Produce json
// @Param name body types.PostTaskHandlerRequest true "Task name"
// @Success 200 {task} types.PostTaskHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Router /task [post]
func (s *Task) postHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatePostTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	newTask := models.Task{Id: uuid.New().String(), Name: req.TaskName, Status: "in_progress", Result:""}

	go editing.Edit(s.service, &newTask)

	err = s.service.Post(newTask)
	types.ProcessError(w, err, &types.PostTaskHandlerResponse{ID: &newTask.Id})
}


// WithtaskHandlers registers task-related HTTP handlers.
func (s *Task) WithTaskHandlers(r chi.Router) {
	r.Route("/task", func(r chi.Router) {
		r.Get("/status/{id}", s.getStatusHandler)
		r.Get("/result/{id}", s.getResultHandler)
		r.Post("/", s.postHandler)
		//r.Put("/", s.putHandler)
		//r.Delete("/", s.deleteHandler)
	})
}