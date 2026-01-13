package mocks

import (
	"code_processor/http_server/models"
	"code_processor/http_server/repository"

	"github.com/google/uuid"
)

type SessionRepo struct {
	Sessions map[string]models.Session
}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{
		Sessions: make(map[string]models.Session),
	}
}

func (sr *SessionRepo) Get(key uuid.UUID) (*models.Session, error) {
	session, exists := sr.Sessions[key.String()]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return &session, nil
}

func (sr *SessionRepo) Set(session models.Session) error {
	sr.Sessions[session.SessionId.String()] = session
	return nil
}

func (sr *SessionRepo) Delete(key uuid.UUID) error {
	delete(sr.Sessions, key.String())
	return nil
}
