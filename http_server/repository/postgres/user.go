package postgres

import (
	"code_processor/http_server/models"
	"code_processor/http_server/repository"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(connStr string) (*UserStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("connecting to db: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging db: %s", err)
	}

	return &UserStorage{db: db}, nil
}

func (us *UserStorage) GetByLogin(login string) (*models.User, error) {
	var user models.User

	err := us.db.QueryRow(`
		SELECT user_id, user_login, user_password
		FROM users
		WHERE user_login = $1`, login).Scan(
		&user.Id,
		&user.Login,
		&user.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound{Item: "user"}
		}
		return nil, fmt.Errorf("querying user by login: %w", err)
	}
	return &user, nil
}

func (us *UserStorage) GetById(key uuid.UUID) (*models.User, error) {
	var user models.User

	err := us.db.QueryRow(`
		SELECT user_id, user_login, user_password
		FROM users
		WHERE user_id = $1`, key.String()).Scan(
		&user.Id,
		&user.Login,
		&user.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound{Item: "user"}
		}
		return nil, fmt.Errorf("querying user by id: %w", err)
	}
	return &user, nil
}

func (us *UserStorage) Put(user models.User) error {
	result, err := us.db.Exec(`
		UPDATE users
		SET user_login = $1,
		    user_password = $2
		WHERE user_id = $3`,
		user.Login,
		user.Password,
		user.Id,
	)

	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return repository.ErrNotFound{Item: "user"}
	}

	return nil
}

func (us *UserStorage) Post(user models.User) error {
	_, err := us.db.Exec(`
		INSERT INTO users (user_id, user_login, user_password)
		VALUES ($1, $2, $3)`,
		user.Id,
		user.Login,
		user.Password,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Constraint == "uq_user_login" {
				return repository.ErrConflict{
					Field: "login",
				}
			}
		}

		return fmt.Errorf("creating user: %w", err)
	}

	return nil
}

func (us *UserStorage) Delete(key uuid.UUID) error {
	result, err := us.db.Exec(`
		DELETE FROM users
		WHERE user_id = $1`, key.String())

	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return repository.ErrNotFound{Item: "user"}
	}

	return nil
}
