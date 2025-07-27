package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func (a *Auth) Login(login, password string) (string, error) {

	uq := `SELECT id, password_hash FROM users WHERE login = $1`

	var userId int64
	var passwordHash string
	err := a.pgPool.QueryRow(context.Background(), uq, login).Scan(&userId, &passwordHash)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("Wrong password")
	}

	tq := `SELECT access_token FROM tokens WHERE user_id = $1 AND access_token_exires_at > NOW()`

	var tokenString string

	err = a.pgPool.QueryRow(context.Background(), tq, userId).Scan(&tokenString)
	if err != nil && err != pgx.ErrNoRows {
		return "", err
	}
	if err != nil && err == pgx.ErrNoRows {
		ts, te, tn, err := GenAccessToken(userId)
		if err != nil {
			return "", err
		}

		tuq := `UPDATE tokens SET access_token = $1, access_token_exires_at = $2, access_token_issued_at = $3 WHERE user_id = $4`

		_, err = a.pgPool.Exec(context.Background(), tuq, ts, te, tn, userId)
		if err != nil {
			return "", err
		}

		tokenString = ts
	}

	return tokenString, nil
}
