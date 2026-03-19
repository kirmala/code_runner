package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/domain"
	"github.com/kirmala/code_runner/http_server/repository"
	"github.com/redis/go-redis/v9"
)

type SessionStorage struct {
	cli *redis.ClusterClient
}

func NewSessionStorage(cli *redis.ClusterClient) (*SessionStorage) {
	return &SessionStorage{cli: cli}
}

func (rs *SessionStorage) Set(ctx context.Context, session domain.Session) error {
	key := fmt.Sprintf("session:%s", session.SessionId.String())
	err := rs.cli.Set(ctx, key, session.UserId.String(), 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("setting session: %w", err)
	}
	return nil
}

func (rs *SessionStorage) Get(ctx context.Context, key uuid.UUID) (*domain.Session, error) {
	fullkey := fmt.Sprintf("session:%s", key.String())
	userIdstr, err := rs.cli.Get(ctx, fullkey).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, repository.ErrNotFound{Item: "session"}
		}
		return nil, fmt.Errorf("getting session: %w", err)
	}

	userId, err := uuid.Parse(userIdstr)
	if err != nil {
		return nil, fmt.Errorf("parsing user id: %w", err)
	}

	return &domain.Session{SessionId: key, UserId: userId}, nil
}

func (rs *SessionStorage) Delete(ctx context.Context, key uuid.UUID) error {
	fullkey := fmt.Sprintf("session:%s", key.String())
	err := rs.cli.Del(ctx, fullkey).Err()
	if err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}
	return nil
}
