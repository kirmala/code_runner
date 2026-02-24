package postgres

import (
	"code_processor/consumer/repository"
	"code_processor/http_server/models"
	"context"
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

func (ps *TaskStorage) Put(ctx context.Context, task models.Task) error {
	result, err := ps.db.ExecContext(ctx, `
		UPDATE tasks 
		SET task_code = $1, 
		    task_translator = $2,
			task_result = $3,
		    task_status = $4
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
