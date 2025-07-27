package handlers

import (
	"billing/internal/psql"
)

type Handler struct {
	pgPool *psql.Pool
}

func New(pgPool *psql.Pool) *Handler {
	return &Handler{pgPool: pgPool}
}
