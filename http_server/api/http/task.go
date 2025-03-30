package http

import (
	"code_processor/http_server/api/http/types"
	"code_processor/http_server/models"
	"code_processor/http_server/usecases"
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Task struct {
	service usecases.Task
}

func NewTaskHandler(service usecases.Task) *Task {
	return &Task{service: service}
}

// @Summary Get a task status
// @Description Get a task status by its id
// @Tags task
// @Accept  json
// @Produce json
// @Param id path string true "Id of the task"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 200 {task} types.GetTaskStatusHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Task not found"
// @Router /task/status/{id} [get]
func (s *Task) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	authToken, err := types.GetAuthToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	_, err = s.service.GetUserId(*authToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	req, err := types.CreateGetTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	taskStatus, err := s.service.GetStatus(req.Id)
	types.ProcessError(w, err, &types.GetTaskStatusHandlerResponse{Status: taskStatus}, 200)
}

// @Summary Get a task result
// @Description Get a task result by its id
// @Tags task
// @Accept  json
// @Produce json
// @Param id path string true "Id of the task"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 200 {task} types.GetTaskResultHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Task not found"
// @Router /task/result/{id} [get]
func (s *Task) getResultHandler(w http.ResponseWriter, r *http.Request) {
	authToken, err := types.GetAuthToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	_, err = s.service.GetUserId(*authToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	req, err := types.CreateGetTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	taskResult, err := s.service.GetResult(req.Id)
	types.ProcessError(w, err, &types.GetTaskResultHandlerResponse{Result: taskResult}, 200)
}

// @Summary Create an task
// @Description Create a new task with the specified key and value
// @Tags task
// @Accept  json
// @Produce json
// @Param translation_data body types.PostTaskHandlerRequest true "task code and translator"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 201 {task} types.PostTaskHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Router /task [post]
func (s *Task) postHandler(w http.ResponseWriter, r *http.Request) {
	authToken, err := types.GetAuthToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	_, err = s.service.GetUserId(*authToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	req, err := types.CreatePostTaskHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	newTask := models.Task{Id: uuid.New().String(), Code: req.TaskCode, Translator: req.TaskTranslator, Status: "in_progress", Result: "progres..."}

	err = s.service.Post(newTask)
	types.ProcessError(w, err, &types.PostTaskHandlerResponse{ID: &newTask.Id}, 201)
}

func (s *Task) commitHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "error while decoding json", http.StatusInternalServerError)
		return
	}
	err := s.service.Put(task)
	if err != nil {
		http.Error(w, "error putting task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
}

// WithtaskHandlers registers task-related HTTP handlers.
func (s *Task) WithTaskHandlers(r chi.Router) {
	r.Route("/task", func(r chi.Router) {
		r.Get("/status/{id}", s.getStatusHandler)
		r.Get("/result/{id}", s.getResultHandler)
		r.Post("/commit", s.commitHandler)
		r.Post("/", s.postHandler)
	})
}
