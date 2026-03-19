package repository

import (
	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/internal/domain"
)

type User interface {
	GetByLogin(login string) (*domain.User, error)
	GetById(id uuid.UUID) (*domain.User, error)
	Put(domain.User) error
	Post(domain.User) error
	Delete(key uuid.UUID) error
}
