package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateConnection(ctx context.Context, connStr string) (*pgx.Conn, error) {
	return pgx.Connect(ctx, connStr)
}
