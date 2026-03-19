package postgres

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/kirmala/code_runner/http_server/domain"
	"github.com/kirmala/code_runner/http_server/repository"
)

type TaskStorage struct {
	db *sql.DB
}

func NewTaskStorage(connStr string) (*TaskStorage, error) {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("connecting to db: %s", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging db: %s", err)
	}

	return &TaskStorage{db: db}, nil
}

func (ps *TaskStorage) Get(key uuid.UUID) (*domain.Task, error) {
	var task domain.Task

	err := ps.db.QueryRow(`
		SELECT task_id, task_code, task_translator, task_status, task_result 
		FROM tasks 
		WHERE task_id = $1`, key.String()).Scan(
		&task.Id,
		&task.Code,
		&task.Translator,
		&task.Status,
		&task.Result,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound{Item: "task"}
		}
		return nil, fmt.Errorf("querying task: %w", err)
	}
	return &task, nil
}

func (ps *TaskStorage) Put(task domain.Task) error {
	result, err := ps.db.Exec(`
		UPDATE tasks 
		SET task_code = $1, 
		    task_translator = $2,
			task_result = $4,
		    task_status = $3
		WHERE task_id = $5`,
		task.Code,
		task.Translator,
		task.Result,
		task.Status,
		task.Id,
	)

	if err != nil {
		return fmt.Errorf("updating task: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return repository.ErrNotFound{Item: "task"}
	}

	return nil
}

func (ps *TaskStorage) Post(task domain.Task) error {
	_, err := ps.db.Exec(`
		INSERT INTO tasks (task_id, task_code, task_translator, task_result, task_status)
		VALUES ($1, $2, $3, $4, $5)`,
		task.Id,
		task.Code,
		task.Translator,
		task.Result,
		task.Status,
	)

	if err != nil {
		return fmt.Errorf("creating task: %w", err)
	}

	return nil
}

func (ps *TaskStorage) Delete(key uuid.UUID) error {
	result, err := ps.db.Exec(`
		DELETE FROM tasks 
		WHERE task_id = $1`, key)

	if err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return repository.ErrNotFound{Item: "task"}
	}

	return nil
}
