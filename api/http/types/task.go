package types

import (
	"github.com/go-chi/chi/v5"
	"encoding/json"
	"fmt"
	"net/http"
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

// type PutObjectHandlerRequest struct {
// 	domain.Object
// }

// func CreatePutObjectHandlerRequest(r *http.Request) (*PutObjectHandlerRequest, error) {
// 	var req PutObjectHandlerRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		return nil, fmt.Errorf("error while decoding json: %v", err)
// 	}
// 	return &req, nil
// }

// type DeleteObjectHandlerRequest struct {
// 	Key string `json:"key"`
// }

// func CreateDeleteObjectHandlerRequest(r *http.Request) (*DeleteObjectHandlerRequest, error) {
// 	key := r.URL.Query().Get("key")
// 	if key == "" {
// 		return nil, fmt.Errorf("missing key")
// 	}
// 	return &DeleteObjectHandlerRequest{Key: key}, nil
// }

type PostTaskHandlerRequest struct {
	TaskName string `json:"task_name"`
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