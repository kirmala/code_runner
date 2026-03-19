package mocks

import (
	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/domain"
	"github.com/kirmala/code_runner/http_server/repository"
)

type SessionRepo struct {
	Sessions map[string]domain.Session
}

func NewSessionRepo() *SessionRepo {
	return &SessionRepo{
		Sessions: make(map[string]domain.Session),
	}
}

func (sr *SessionRepo) Get(key uuid.UUID) (*domain.Session, error) {
	session, exists := sr.Sessions[key.String()]
	if !exists {
		return nil, repository.ErrNotFound{}
	}
	return &session, nil
}

func (sr *SessionRepo) Set(session domain.Session) error {
	sr.Sessions[session.SessionId.String()] = session
	return nil
}

func (sr *SessionRepo) Delete(key uuid.UUID) error {
	delete(sr.Sessions, key.String())
	return nil
}
