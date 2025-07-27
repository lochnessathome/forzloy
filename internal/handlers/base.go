package handlers

import (
	"billing/internal/psql"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	pgPool *psql.Pool
}

func New(pgPool *psql.Pool) *Handler {
	return &Handler{pgPool: pgPool}
}

func ParseJWTSubject(c echo.Context) string {
	u := c.Get("user")
	claims := u.(*jwt.Token).Claims.(jwt.MapClaims)
	subject := claims["sub"].(string)

	return subject
}
