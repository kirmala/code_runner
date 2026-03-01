package models

import "github.com/google/uuid"

type Task struct {
	Id         uuid.UUID  `json:"id"`
	Code       string     `json:"code"`
	Translator Translator `json:"translator"`
	Status     status     `json:"status"`
	Result     string     `json:"result"`
	//UserId     string `json:"user_id"`
}
