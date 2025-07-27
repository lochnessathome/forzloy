package reports

import (
	"billing/internal/mng"
	"billing/internal/psql"
)

type Reports struct {
	pgPool     *psql.Pool
	mnDatabase *mng.Database
}

func New(pgPool *psql.Pool, mnDatabase *mng.Database) *Reports {
	return &Reports{pgPool: pgPool, mnDatabase: mnDatabase}
}
