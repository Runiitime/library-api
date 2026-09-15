package config

import (
	"fmt"

	"github.com/cloudresty/go-env"
)

func GetConfig() string {
	if err := env.Load(); err != nil {
		panic(err)
	}
	return getConnectionString()
}

func getConnectionString() string {
	dbName := env.Get("DB_NAME", "postgres")
	dbUser := env.Get("DB_USER", "postgres")
	dbPassword := env.Get("DB_PASSWORD", "postgres")
	dbHost := env.Get("DB_HOST", "localhost")
	dbPort := env.Get("DB_PORT", "5432")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
}
