package models

import "github.com/google/uuid"

type Task struct {
	Id         uuid.UUID `json:"id"`
	Code       string `json:"code"`
	Translator string `json:"translator"`
	Status     string `json:"status"`
	Result     string `json:"result"`
	//UserId     string `json:"user_id"`
}
