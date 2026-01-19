package httpx

import (
	"code_processor/http_server/api"
	"code_processor/http_server/api/dto"
	"code_processor/http_server/models"
	"code_processor/http_server/usecases"
	"fmt"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateGetTaskHandlerRequest(r *http.Request) *dto.GetTaskHandlerRequest {
	id := chi.URLParam(r, "id")
	return &dto.GetTaskHandlerRequest{Id: id}
}

func CreatePostTaskHandlerRequest(r *http.Request) (*dto.PostTaskHandlerRequest, error) {
	var req dto.PostTaskHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, api.ErrBadRequest{}
	}
	return &req, nil
}

func GetAuthToken(r *http.Request) (uuid.UUID, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return uuid.Nil, fmt.Errorf("invalid Authorization token format need to include Bearer")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := uuid.Parse(tokenStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid Authorization token format: %v", err)
	}

	return token, nil
}

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
	authToken, err := GetAuthToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	_, err = s.service.GetUserId(authToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	req := CreateGetTaskHandlerRequest(r)

	id, err := uuid.Parse(req.Id)
	if err != nil {
		WriteResponse(w, api.ErrBadRequest{Field: "id", Err: err.Error()}, nil)
		return
	}

	taskStatus, err := s.service.GetStatus(id)
	WriteResponse(w, err, &dto.GetTaskStatusHandlerResponse{Status: taskStatus})
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
	authToken, err := GetAuthToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	_, err = s.service.GetUserId(authToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	req := CreateGetTaskHandlerRequest(r)

	id, err := uuid.Parse(req.Id)

	if err != nil {
		WriteResponse(w, api.ErrBadRequest{Field: "id", Err: err.Error()}, nil)
		return
	}

	taskResult, err := s.service.GetResult(id)
	WriteResponse(w, err, dto.GetTaskResultHandlerResponse{Result: taskResult})
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
	authToken, err := GetAuthToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	_, err = s.service.GetUserId(authToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	req, err := CreatePostTaskHandlerRequest(r)
	if err != nil {
		WriteResponse(w, err, nil)
		return
	}

	taskTranslator, err := models.ParseTranslator(req.TaskTranslator)

	if err != nil {
		WriteResponse(w, api.ErrBadRequest{Field: "task_translator", Err: err.Error()}, nil)
		return
	}

	newTask := models.Task{Id: uuid.New(), Code: req.TaskCode, Translator: taskTranslator, Status: models.StatusInProgress, Result: "progres..."}

	err = s.service.Post(newTask)
	WriteResponse(w, err, dto.PostTaskHandlerResponse{ID: newTask.Id.String()})
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
