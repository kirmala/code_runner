package postgres

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// RunMigrations executes all pending migrations in the "migrations" directory using the provided connection string to connect to the database.
func RunMigrations(connStr string) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return fmt.Errorf("connecting to db: %s", err)
	}
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("pinging db: %s", err)
	}

	return goose.Up(db, "migrations")
}
