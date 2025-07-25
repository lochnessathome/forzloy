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

func New(ctx context.Context) (*Pool, error) {
	dbUrl := os.Getenv("DATABASE_URL")
        if dbUrl == "" {
                return nil, fmt.Errorf("Empty DATABASE_URL environment variable")
        }

        pool, err := pgxpool.New(ctx, dbUrl)
        if err != nil {
		return nil, err
        }

	return &Pool{pool}, nil
}

