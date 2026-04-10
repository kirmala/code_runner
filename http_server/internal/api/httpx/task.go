package httpx

import (
	"net/http"

	"encoding/json"

	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/internal/api"
	"github.com/kirmala/code_runner/http_server/internal/api/dto"
	"github.com/kirmala/code_runner/http_server/internal/api/httpx/middleware"
	"github.com/kirmala/code_runner/http_server/internal/domain"
	"github.com/kirmala/code_runner/http_server/internal/service"
	"github.com/labstack/echo/v5"
)

func CreateGetTaskHandlerRequest(c *echo.Context) *dto.GetTaskHandlerRequest {
	id := c.Param("id")

	return &dto.GetTaskHandlerRequest{Id: id}
}

func CreatePostTaskHandlerRequest(r *http.Request) (*dto.PostTaskHandlerRequest, error) {
	var req dto.PostTaskHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, api.ErrBadRequest{}
	}
	return &req, nil
}

type Task struct {
	service service.Task
	auth service.Authenticator
}

func NewTaskHandler(service service.Task, auth service.Authenticator) *Task {
	return &Task{service: service, auth: auth}
}

// @Summary Get a task status
// @Description Get a task status by its id
// @Tags task
// @Accept  json
// @Produce json
// @Param id path string true "Id of the task"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 200 {object} dto.GetTaskStatusHandlerResponse
// @Failure 400 {object} dto.Error "Bad request"
// @Failure 404 {object} dto.Error "Task not found"
// @Failure 401 {object} dto.Error "Unauthorized"
// @Router /task/status/{id} [get]
func (s *Task) getStatusHandler(c *echo.Context) error {
	_, ok := c.Get(middleware.UserIdKey).(uuid.UUID)
	if !ok {
		panic("getStatusHandler called with no user id in context")
	}

	req := CreateGetTaskHandlerRequest(c)

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return api.ErrBadRequest{Field: "id", Err: err.Error()}
	}

	taskStatus, err := s.service.GetStatus(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &dto.GetTaskStatusHandlerResponse{Status: taskStatus})
}

// @Summary Get a task result
// @Description Get a task result by its id
// @Tags task
// @Accept  json
// @Produce json
// @Param id path string true "Id of the task"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 200 {object} dto.GetTaskResultHandlerResponse
// @Failure 400 {object} dto.Error "Bad request"
// @Failure 404 {object} dto.Error "Task not found"
// @Failure 401 {object} dto.Error "Unauthorized"
// @Router /task/result/{id} [get]
func (s *Task) getResultHandler(c *echo.Context) error {
	_, ok := c.Get(middleware.UserIdKey).(uuid.UUID)
	if !ok {
		panic("getResultHandler called with no user id in context")
	}

	req := CreateGetTaskHandlerRequest(c)

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return api.ErrBadRequest{Field: "id", Err: err.Error()}
	}

	
	taskResult, err := s.service.GetResult(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &dto.GetTaskResultHandlerResponse{Result: taskResult})
}

// @Summary Create an task
// @Description Create a new task with the specified key and value
// @Tags task
// @Accept  json
// @Produce json
// @Param translation_data body dto.PostTaskHandlerRequest true "task code and translator"
// @Param Authorization header string true "Bearer token for authentication"
// @Success 201 {object} dto.PostTaskHandlerResponse
// @Failure 400 {object} dto.Error "Bad request"
// @Failure 401 {object} dto.Error "Unauthorized"
// @Router /task [post]
func (s *Task) postHandler(c *echo.Context) error {
	req, err := CreatePostTaskHandlerRequest(c.Request())
	if err != nil {
		return err
	}

	taskTranslator, err := domain.ParseTranslator(req.TaskTranslator)

	if err != nil {
		return api.ErrBadRequest{Field: "task_translator", Err: err.Error()}
	}

	newTask := domain.Task{Id: uuid.New(), Code: req.TaskCode, Translator: taskTranslator, Status: domain.StatusInProgress, Result: "progres..."}

	err = s.service.Post(c.Request().Context(), newTask)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, &dto.PostTaskHandlerResponse{ID: newTask.Id.String()})
}

// WithtaskHandlers registers task-related HTTP handlers.
func (s *Task) WithTaskHandlers(g *echo.Group) {
	taskg := g.Group("/task", middleware.Auth{Authenticator: s.auth}.Auth)

	taskg.GET("/status/:id", s.getStatusHandler)
	taskg.GET("/result/:id", s.getResultHandler)
	taskg.POST("", s.postHandler)
}
