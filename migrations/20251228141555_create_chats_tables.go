package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateChatTables, downCreateChatTables)
}

func upCreateChatTables(ctx context.Context, tx *sql.Tx) error {
	query := `
        CREATE TABLE chats (
            id BIGSERIAL PRIMARY KEY,
            type VARCHAR(20) NOT NULL,
            name VARCHAR(255),
            created_at TIMESTAMP NOT NULL DEFAULT NOW(),
            CHECK (type IN ('private', 'group'))
        );

        CREATE TABLE chat_members (
            chat_id BIGINT NOT NULL,
            user_id BIGINT NOT NULL,
            joined_at TIMESTAMP NOT NULL DEFAULT NOW(),

            PRIMARY KEY (chat_id, user_id),

            FOREIGN KEY (chat_id)
                REFERENCES chats(id)
                ON DELETE CASCADE,

            FOREIGN KEY (user_id)
                REFERENCES users(id)
                ON DELETE CASCADE
        );
    `
	_, err := tx.ExecContext(ctx, query)
	return err
}

func downCreateChatTables(ctx context.Context, tx *sql.Tx) error {
	query := `
        DROP TABLE chat_members;
        DROP TABLE chats;
    `
	_, err := tx.ExecContext(ctx, query)
	return err
}
