package basic

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/domain"
	"github.com/kirmala/code_runner/http_server/repository"
	"github.com/kirmala/code_runner/http_server/service"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	userRepo    repository.User
	sessionRepo repository.Session
}

func NewUser(userRepo repository.User, sessionRepo repository.Session) *User {
	return &User{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (rs *User) Get(key uuid.UUID) (*domain.User, error) {
	user, err := rs.userRepo.GetById(key)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (rs *User) Register(user domain.User) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	return rs.userRepo.Post(domain.User{Id: user.Id, Login: user.Login, Password: string(hashedPassword)})
}

func (rs *User) Login(ctx context.Context, login string, password string) (uuid.UUID, error) {
	user, err := rs.userRepo.GetByLogin(login)
	if err != nil {
		var target repository.ErrNotFound
		if errors.As(err, &target) {
			return uuid.Nil, service.ErrUnauthenticated{Msg: target.Error()}
		}
		return uuid.Nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return uuid.Nil, service.ErrUnauthenticated{Msg: "password is incorrect"}
	}
	sessionId := uuid.New()
	err = rs.sessionRepo.Set(ctx, domain.Session{UserId: user.Id, SessionId: sessionId})
	if err != nil {
		return uuid.Nil, err
	}
	return sessionId, nil
}

func (rs *User) Delete(key uuid.UUID) error {
	return rs.userRepo.Delete(key)
}
