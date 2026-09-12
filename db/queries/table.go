package queries

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateTable(conn *pgx.Conn, ctx context.Context, table string) error {
	q := `
	CREATE TABLE IF NOT EXISTS ` + table + ` (
	id serial PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	review VARCHAR(1000),
	author VARCHAR(255) NOT NULL,
	pages INTEGER NOT NULL,
	published VARCHAR(4) NOT NULL,
	completed BOOLEAN NOT NULL,
	created_at TIMESTAMP NOT NULL,
	completed_at TIMESTAMP,

	UNIQUE (title)
	);`

	_, err := conn.Exec(ctx, q)
	return err
}
