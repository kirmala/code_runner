package service

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"photo_editor/models"
	"photo_editor/repository"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	userRepo repository.User
	sessionRepo repository.Session
}

func NewUser(userRepo repository.User, sessionRepo repository.Session) *User {
	return &User{
		userRepo: userRepo,
		sessionRepo: sessionRepo,
	}
}

func sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}


func (rs *User) Get(key string) (*models.User, error) {
	user, err := rs.userRepo.Get(key)
	if (err != nil) {
		return nil, err
	}
	return user, err
}

func (rs *User) PostRegister(user models.User) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	return rs.userRepo.Post(models.User{Id: user.Id, Login: user.Login, Password: string(hashedPassword)})
}

func (rs *User) PostLogin(login string, password string) (*string, error) {
	user, err := rs.userRepo.Get(login)
	if (err != nil) {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if (err != nil) {
		return nil, err
	}
	sessionId := sessionId()
	err = rs.sessionRepo.Post(models.Session{UserId: user.Id, SessionId: sessionId})
	if (err != nil) {
		return nil, err
	}
	return &sessionId, nil
}

func (rs *User) Delete(key string) error {
	return rs.userRepo.Delete(key)
}


