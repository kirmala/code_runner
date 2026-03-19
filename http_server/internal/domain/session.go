package domain

import "github.com/google/uuid"

type Session struct {
	UserId    uuid.UUID `json:"user_id"`
	SessionId uuid.UUID `json:"session_id"`
}
