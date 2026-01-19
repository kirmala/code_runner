package redis

import (
	"code_processor/http_server/models"
	"code_processor/http_server/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionStorage struct {
	db *redis.Client
}

func NewSessionStorage(addr string, password string) (*SessionStorage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	var ctx = context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("pinging redis: %w", err)
	}
	return &SessionStorage{db: rdb}, nil
}

func (rs *SessionStorage) Set(session models.Session) error {
	var ctx = context.Background()
	err := rs.db.Set(ctx, session.SessionId.String(), session.UserId.String(), 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("setting session: %w", err)
	}
	return nil
}

func (rs *SessionStorage) Get(key uuid.UUID) (*models.Session, error) {
	var ctx = context.Background()
	userIdstr, err := rs.db.Get(ctx, key.String()).Result()

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

	return &models.Session{SessionId: key, UserId: userId}, nil
}

func (rs *SessionStorage) Delete(key uuid.UUID) error {
	var ctx = context.Background()
	err := rs.db.Del(ctx, key.String()).Err()
	if err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}
	return nil
}
