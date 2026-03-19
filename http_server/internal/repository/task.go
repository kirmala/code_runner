package repository

import (
	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/internal/domain"
)

type Task interface {
	Get(key uuid.UUID) (*domain.Task, error)
	Put(domain.Task) error
	Post(domain.Task) error
	Delete(key uuid.UUID) error
}
