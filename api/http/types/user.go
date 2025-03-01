package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetAuthToken (r *http.Request) (*string) {
	authHeader := r.Header.Get("Authorization")
	return &authHeader
}

type PostUserRegisterHandlerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreatePostUserRegisterHandlerRequest(r *http.Request) (*PostUserRegisterHandlerRequest, error) {
	var req PostUserRegisterHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

type PostUserLoginHandlerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PostUserLoginHandlerResponse struct {
	Token *string `json:"token"`
}

func CreatePostUserLoginHandlerRequest(r *http.Request) (*PostUserLoginHandlerRequest, error) {
	var req PostUserLoginHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

