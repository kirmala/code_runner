package ram_storage

import (
	"code_runner/models"
	"code_runner/repository"
)

type User struct {
	data map[string]models.User
}

func NewUser() *User {
	return &User{
		data: make(map[string]models.User),
	}
}

func (rs *User) Get(key string) (*models.User, error) {
	value, exists := rs.data[key]
	if !exists {
		return nil, repository.NotFound
	}
	return &value, nil
}

func (rs *User) Put(user models.User) error {
	rs.data[user.Login] = user
	return nil
}

func (rs *User) Post(user models.User) error {
	if _, exists := rs.data[user.Login]; exists {
		return repository.AlreadyExists
	}
	rs.data[user.Login] = user
	return nil
}

func (rs *User) Delete(key string) error {
	if _, exists := rs.data[key]; !exists {
		return repository.NotFound
	}
	delete(rs.data, key)
	return nil
}
