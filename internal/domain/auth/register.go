package auth

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const accessTokenLifetime = 60 * 24 * time.Hour

func (a *Auth) Register(login, password string) (string, error) {

	tx, err := a.pgPool.Begin(context.Background())
	if err != nil {
		return "", err
	}

	defer tx.Rollback(context.Background())

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	uq := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`

	var id int64
	err = tx.QueryRow(context.Background(), uq, login, passwordHash).Scan(&id)
	if err != nil {
		return "", err
	}

	tn := time.Now()
	te := tn.Add(accessTokenLifetime)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(te),
		IssuedAt:  jwt.NewNumericDate(tn),
		NotBefore: jwt.NewNumericDate(tn),
		Subject:   strconv.FormatInt(id, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	tq := `INSERT INTO tokens (user_id, access_token, access_token_exires_at, access_token_issued_at) VALUES ($1, $2, $3, $4)`

	_, err = tx.Exec(context.Background(), tq, id, tokenString, te, tn)
	if err != nil {
		return "", err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
