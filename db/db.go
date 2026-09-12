package db

import (
	"context"
	"fmt"

	"github.com/cloudresty/go-env"
	"github.com/jackc/pgx/v5"
)

func CreateConnection(ctx context.Context) (*pgx.Conn, error) {
	if err := env.Load(); err != nil {
		panic(err)
	}

	connStr := getConnectionString()
	fmt.Println(connStr)

	return pgx.Connect(ctx, connStr)
}

func getConnectionString() string {
	dbName := env.Get("DB_NAME", "postgres")
	dbUser := env.Get("DB_USER", "postgres")
	dbPassword := env.Get("DB_PASSWORD", "postgres")
	dbUrl := env.Get("DB_URL", "postgres")
	dbPort := env.Get("DB_PORT", "5432")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbUrl, dbPort, dbName)
}
