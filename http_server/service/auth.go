package service

import (
	"github.com/google/uuid"
)

type Authenticator interface {
	Authenticate(token string) (uuid.UUID, error)
}