package types

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type GetTaskHandlerRequest struct {
	Id string `json:"id"`
}

type GetTaskStatusHandlerResponse struct {
	Status *string `json:"status"`
}

type GetTaskResultHandlerResponse struct {
	Result *string `json:"result"`
}

func CreateGetTaskHandlerRequest(r *http.Request) (*GetTaskHandlerRequest, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return nil, fmt.Errorf("missing id")
	}
	return &GetTaskHandlerRequest{Id: id}, nil
}

type PostTaskHandlerRequest struct {
	TaskCode       string `json:"task_code"`
	TaskTranslator string `json:"task_translator"`
}

type PostTaskHandlerResponse struct {
	ID *string `json:"id"`
}

func CreatePostTaskHandlerRequest(r *http.Request) (*PostTaskHandlerRequest, error) {
	var req PostTaskHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}
