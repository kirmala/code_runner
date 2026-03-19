package ram_storage

// import (
// 	"http_server/domain"
// 	"http_server/repository"

// 	"github.com/google/uuid"
// )

// type User struct {
// 	data map[uuid.UUID]domain.User
// }

// func NewUser() *User {
// 	return &User{
// 		data: make(map[uuid.UUID]domain.User),
// 	}
// }

// func (u *User) GetByLogin(login string) (*domain.User, error) {
// 	value, exists := rs.data[login]
// 	if !exists {
// 		return nil, repository.ErrNotFound{Item: "user"}
// 	}
// 	return &value.Status, nil
// }

// func (u *User) GetResult(key uuid.UUID) (string, error) {
// 	value, exists := rs.data[key]
// 	if !exists {
// 		return nil, repository.ErrNotFound{Item: "user"}
// 	}
// 	return &value.Result, nil
// }

// func (rs *User) Get(key uuid.UUID) (*domain.User, error) {
// 	value, exists := rs.data[key]
// 	if !exists {
// 		return nil, repository.ErrNotFound{Item: "user"}
// 	}
// 	return &value, nil
// }

// func (rs *User) Put(user domain.User) error {
// 	rs.data[user.Id] = user
// 	return nil
// }

// func (rs *User) Post(user domain.User) error {
// 	if _, exists := rs.data[user.Login]; exists {
// 		return repository.ErrAlreadyExists
// 	}
// 	rs.data[user.Login] = user
// 	return nil
// }

// func (rs *User) Delete(key string) error {
// 	if _, exists := rs.data[key]; !exists {
// 		return repository.ErrNotFound
// 	}
// 	delete(rs.data, key)
// 	return nil
// }
