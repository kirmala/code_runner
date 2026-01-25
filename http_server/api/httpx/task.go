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
// @Success 200 {object} dto.GetTaskStatusHandlerResponse
// @Failure 400 {object} HTTPError "Bad request"
// @Failure 404 {object} HTTPError "Task not found"
// @Failure 401 {object} HTTPError "Unauthorized"
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
		WriteError(w, api.ErrBadRequest{Field: "id", Err: err.Error()})
		return
	}

	taskStatus, err := s.service.GetStatus(id)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteSuccess(w, &dto.GetTaskStatusHandlerResponse{Status: taskStatus}, http.StatusOK)
}

// @Summary Get a task result
// @Description Get a task result by its id
// @Tags task
// @Accept  json
// @Produce json
// @Param id path string true "Id of the task"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 200 {task} dto.GetTaskResultHandlerResponse
// @Failure 400 {object} HTTPError "Bad request"
// @Failure 404 {object} HTTPError "Task not found"
// @Failure 401 {object} HTTPError "Unauthorized"
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
		WriteError(w, api.ErrBadRequest{Field: "id", Err: err.Error()})
		return
	}

	taskResult, err := s.service.GetResult(id)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteSuccess(w, dto.GetTaskResultHandlerResponse{Result: taskResult}, http.StatusOK)
}

// @Summary Create an task
// @Description Create a new task with the specified key and value
// @Tags task
// @Accept  json
// @Produce json
// @Param translation_data body types.PostTaskHandlerRequest true "task code and translator"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 201 {task} types.PostTaskHandlerResponse
// @Failure 400 {object} HTTPError "Bad request"
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
		WriteError(w, err)
		return
	}

	taskTranslator, err := models.ParseTranslator(req.TaskTranslator)

	if err != nil {
		WriteError(w, api.ErrBadRequest{Field: "task_translator", Err: err.Error()})
		return
	}

	newTask := models.Task{Id: uuid.New(), Code: req.TaskCode, Translator: taskTranslator, Status: models.StatusInProgress, Result: "progres..."}

	err = s.service.Post(newTask)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteSuccess(w, dto.PostTaskHandlerResponse{ID: newTask.Id.String()}, http.StatusCreated)
}

// WithtaskHandlers registers task-related HTTP handlers.
func (s *Task) WithTaskHandlers(r chi.Router) {
	r.Route("/task", func(r chi.Router) {
		r.Get("/status/{id}", s.getStatusHandler)
		r.Get("/result/{id}", s.getResultHandler)
		r.Post("/", s.postHandler)
	})
}
