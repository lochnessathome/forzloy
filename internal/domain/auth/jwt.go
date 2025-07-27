package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const accessTokenLifetime = 60 * 24 * time.Hour

func GenAccessToken(userId int64) (string, time.Time, time.Time, error) {

	tn := time.Now()
	te := tn.Add(accessTokenLifetime)

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(te),
		IssuedAt:  jwt.NewNumericDate(tn),
		NotBefore: jwt.NewNumericDate(tn),
		Subject:   strconv.FormatInt(userId, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, time.Time{}, err
	}

	return tokenString, te, tn, nil
}
