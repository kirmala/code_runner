package ram_storage

import (
	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/internal/domain"
	"github.com/kirmala/code_runner/http_server/internal/repository"
)

type Session struct {
	data map[uuid.UUID]domain.Session
}

func NewSession() *Session {
	return &Session{
		data: make(map[uuid.UUID]domain.Session),
	}
}

func (rs Session) Get(key uuid.UUID) (*domain.Session, error) {
	value, exists := rs.data[key]
	if !exists {
		return nil, repository.ErrNotFound{Item: "session"}
	}
	return &value, nil
}

func (rs *Session) Set(session domain.Session) error {
	rs.data[session.SessionId] = session
	return nil
}

func (rs *Session) Delete(key uuid.UUID) error {
	delete(rs.data, key)
	return nil
}
