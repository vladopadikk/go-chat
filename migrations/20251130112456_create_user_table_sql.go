package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUserTableSql, downCreateUserTableSql)
}

func upCreateUserTableSql(ctx context.Context, tx *sql.Tx) error {
	query := `
			CREATE TABLE users (
				id SERIAL PRIMARY KEY,
				email VARCHAR(255) UNIQUE NOT NULL,
				username VARCHAR(100) NOT NULL,
				password_hash TEXT NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT NOW()
			);
	`
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateUserTableSql(ctx context.Context, tx *sql.Tx) error {
	query := `DROP TABLE users;`
	_, err := tx.ExecContext(ctx, query)
	return err
}
