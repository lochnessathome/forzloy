package auth

import (
	"billing/internal/psql"
)

type Auth struct {
	pgPool *psql.Pool
}

func New(pgPool *psql.Pool) *Auth {
	return &Auth{pgPool: pgPool}
}
