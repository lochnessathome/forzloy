package psql

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
}

func New() (*Pool, error) {
	dbUrl := os.Getenv("PSQL_DATABASE_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("Empty PSQL_DATABASE_URL environment variable")
	}

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	return &Pool{pool}, nil
}
