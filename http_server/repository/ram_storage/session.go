package ram_storage

import (
	"code_processor/http_server/models"
	"code_processor/http_server/repository"

	"github.com/google/uuid"
)

type Session struct {
	data map[uuid.UUID]models.Session
}

func NewSession() *Session {
	return &Session{
		data: make(map[uuid.UUID]models.Session),
	}
}

func (rs Session) Get(key uuid.UUID) (*models.Session, error) {
	value, exists := rs.data[key]
	if !exists {
		return nil, repository.ErrNotFound{Item: "session"}
	}
	return &value, nil
}

func (rs *Session) Set(session models.Session) error {
	rs.data[session.SessionId] = session
	return nil
}

func (rs *Session) Delete(key uuid.UUID) error {
	delete(rs.data, key)
	return nil
}
