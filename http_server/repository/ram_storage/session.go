package ram_storage

import (
	"code_processor/http_server/models"
	"code_processor/http_server/repository"
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
		return nil, repository.NotFound
	}
	return &value, nil
}

func (rs *Session) Post(session models.Session) error {
	if _, exists := rs.data[session.SessionId]; exists {
		return repository.AlreadyExists
	}
	rs.data[session.SessionId] = session
	return nil
}

func (rs *Session) Delete(key string) error {
	if _, exists := rs.data[key]; !exists {
		return repository.NotFound
	}
	delete(rs.data, key)
	return nil
}
