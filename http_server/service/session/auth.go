package session

import (
	"context"

	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/repository"
	"github.com/kirmala/code_runner/http_server/service"
)

type Authenticator struct {
	SessionRepo repository.Session
}

func (a Authenticator) Authenticate(ctx context.Context, token string) (uuid.UUID, error) {
	tokenUUID, err  := uuid.Parse(token)

	if err != nil {
		return uuid.Nil, service.ErrUnauthenticated{Msg: "Invalid uuid format"}
	}

	s, err := a.SessionRepo.Get(ctx, tokenUUID)

	if err != nil {
		return uuid.Nil, service.ErrUnauthenticated{Msg: err.Error()}
	}

	uId := s.UserId

	return uId, nil
}