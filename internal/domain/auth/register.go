package auth

import (
	"context"
	//        "fmt"
	//        "os"
	//        "time"

	// "billing/internal/psql"

	"golang.org/x/crypto/bcrypt"
	// "github.com/golang-jwt/jwt/v5"
)

func (a *Auth) Register(login, password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2)`

	_, err = a.pgPool.Exec(context.Background(), query, login, passwordHash)
	if err != nil {
		return "", err
	}

	return "DONE", nil
}
