package ram_storage

import (
	"code_processor/http_server/models"
	"code_processor/http_server/repository"

	"github.com/google/uuid"
)

type Session struct {
	data map[string]models.Session
}

func NewSession() *Session {
	return &Session{
		data: make(map[string]models.Session),
	}
}

func (rs *Session) Get(key string) (*models.Session, error) {
	value, exists := rs.data[key]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return &value, nil
}

func (rs *Session) Post(session models.Session) error {
	if _, exists := rs.data[session.SessionId.String()]; exists {
		return repository.ErrAlreadyExists
	}
	rs.data[session.SessionId.String()] = session
	return nil
}

func (rs *Session) Delete(key uuid.UUID) error {
	if _, exists := rs.data[key.String()]; !exists {
		return repository.ErrNotFound
	}
	delete(rs.data, key.String())
	return nil
}
