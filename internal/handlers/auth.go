package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"billing/internal/domain/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponce struct {
	AccessToken string `json:"access_token"`
}

func AuthCreateToken(c echo.Context) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Now().Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, struct{ Token string }{Token: tokenString})
}

func AuthValidateToken(c echo.Context) error {
	tokenString := c.Param("ts")
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("Bad token")
	}

	return c.JSON(http.StatusOK, claims)
}

func (h *Handler) AuthRegister(c echo.Context) error {
	req := new(RegisterRequest)

	err := c.Bind(req)
	if err != nil {
		return err
	}

	err = c.Validate(req)
	if err != nil {
		return err
	}

	a := auth.New(h.pgPool)

	accessToken, err := a.Register(req.Login, req.Password)
	if err != nil {
		return err
	}

	res := &RegisterResponce{AccessToken: accessToken}

	return c.JSON(http.StatusOK, res)
}
