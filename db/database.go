package db

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func InitDbPool() error {
	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return fmt.Errorf("DATABASE_URL environment variable not set")
	}
	pools, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	pool = pools
	return nil
}

func GetDbPool() (*pgxpool.Pool, error) {
	if pool == nil {
		return nil, errors.New("database connection pool not initialized")
	}
	return pool, nil
}
