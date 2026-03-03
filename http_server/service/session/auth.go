package session

import (
	"code_processor/http_server/repository"
	"code_processor/http_server/service"
	"context"

	"github.com/google/uuid"
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