package mocks

import (
	"code_processor/http_server/models"
	"code_processor/http_server/repository"
)

type SessionRepo struct {
	Sessions map[string]models.Session
}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{
		Sessions: make(map[string]models.Session),
	}
}

func (sr *SessionRepo) Get(key string) (*models.Session, error) {
	session, exists := sr.Sessions[key]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return &session, nil
}

func (sr *SessionRepo) Set(session models.Session) error {
	sr.Sessions[session.SessionId] = session
	return nil
}

func (sr *SessionRepo) Delete(key string) error {
	delete(sr.Sessions, key)
	return nil
}
