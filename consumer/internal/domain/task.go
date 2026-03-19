package domain

import (
	"github.com/google/uuid"
)

type Task struct {
	Id         uuid.UUID  `json:"id"`
	Code       string     `json:"code"`
	Translator Translator `json:"translator"`
	Status     Status     `json:"status"`
	Result     string     `json:"result"`
}