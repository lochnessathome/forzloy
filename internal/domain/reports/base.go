package reports

import (
	"billing/internal/psql"
)

type Reports struct {
	pgPool *psql.Pool
}

func New(pgPool *psql.Pool) *Reports {
	return &Reports{pgPool: pgPool}
}
