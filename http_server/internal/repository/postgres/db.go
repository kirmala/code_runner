package postgres

import (
	"database/sql"
	"fmt"
)

func Connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("db open failed, %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db ping failed, %s", err)
	}

	return db, nil
}
