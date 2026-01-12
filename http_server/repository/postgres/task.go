package postgres

import (
	"code_processor/http_server/models"
	"code_processor/http_server/repository"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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

func (ps *TaskStorage) Get(key string) (*models.Task, error) {
	var task models.Task

	err := ps.db.QueryRow(`
		SELECT task_id, task_code, task_translator, task_status, task_result 
		FROM tasks 
		WHERE task_id = $1`, key).Scan(
		&task.Id,
		&task.Code,
		&task.Translator,
		&task.Status,
		&task.Result,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.NotFound
		}
		return nil, fmt.Errorf("querying task: %w", err)
	}
	return &task, nil
}

func (ps *TaskStorage) Put(task models.Task) error {
	result, err := ps.db.Exec(`
		UPDATE tasks 
		SET task_code = $1, 
		    task_translator = $2,
			task_result = $4,
		    task_status = $3,
		WHERE task_id = $5`,
		task.Code,
		task.Translator,
		task.Result,
		task.Status,
		task.Id,
	)

	if err != nil {
		return fmt.Errorf("pdating task: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no task found with id %s", task.Id)
	}

	return nil
}

func (ps *TaskStorage) Post(task models.Task) error {
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

func (ps *TaskStorage) Delete(key string) error {
	result, err := ps.db.Exec(`
		DELETE FROM tasks 
		WHERE task_id = $1`, key)

	if err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no task found with id %s", key)
	}

	return nil
}
