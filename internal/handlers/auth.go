package handlers

import (
	"net/http"

	"billing/internal/domain/auth"

	"github.com/labstack/echo/v4"
)

type AuthRegisterRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthLoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthRegisterResponce struct {
	AccessToken string `json:"access_token"`
}

type AuthLoginResponce struct {
	AccessToken string `json:"access_token"`
}

/*
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
*/

func (h *Handler) AuthRegister(c echo.Context) error {
	req := new(AuthRegisterRequest)

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

	res := &AuthRegisterResponce{AccessToken: accessToken}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) AuthLogin(c echo.Context) error {
	req := new(AuthLoginRequest)

	err := c.Bind(req)
	if err != nil {
		return err
	}

	err = c.Validate(req)
	if err != nil {
		return err
	}

	a := auth.New(h.pgPool)

	accessToken, err := a.Login(req.Login, req.Password)
	if err != nil {
		return err
	}

	res := &AuthLoginResponce{AccessToken: accessToken}

	return c.JSON(http.StatusOK, res)
}
