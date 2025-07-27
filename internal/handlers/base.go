package handlers

import (
	"billing/internal/mng"
	"billing/internal/psql"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	pgPool     *psql.Pool
	mnDatabase *mng.Database
}

func New(pgPool *psql.Pool, mnDatabase *mng.Database) *Handler {
	return &Handler{pgPool: pgPool, mnDatabase: mnDatabase}
}

func ParseJWTSubject(c echo.Context) string {
	u := c.Get("user")
	claims := u.(*jwt.Token).Claims.(jwt.MapClaims)
	subject := claims["sub"].(string)

	return subject
}
