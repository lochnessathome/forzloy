package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"billing/cmd/migrations"
	"billing/internal/psql"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const httpPort = "8080"

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

        err := migrations.PsqlUp()
        if err != nil {
                e.Logger.Error(err)
                os.Exit(1)
        }

	ctx := context.Background()
	pgPool, err := psql.New(ctx)
	if err != nil {
		e.Logger.Error(err)
		os.Exit(1)
	}

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.GET("/time", func(c echo.Context) error {
		var tm string
		err = pgPool.QueryRow(ctx, "SELECT NOW()::text;").Scan(&tm)
		if err != nil {
			e.Logger.Error(err)
		}

		return c.JSON(http.StatusOK, struct{ Time string }{Time: tm})
	})

	e.GET("/create-token", func(c echo.Context) error {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"foo": "bar",
			"nbf": time.Now().Unix(),
		})

		secret := os.Getenv("JWT_SECRET")
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			e.Logger.Error(err)
		}

		return c.JSON(http.StatusOK, struct{ Token string }{Token: tokenString})
	})

	e.GET("/validate-token/:ts", func(c echo.Context) error {
		tokenString := c.Param("ts")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			secret := os.Getenv("JWT_SECRET")
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil {
			e.Logger.Error(err)

			// нужно останавливать хэндлер
		}



		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
                        e.Logger.Error(fmt.Errorf("Bad token"))
		}

		return c.JSON(http.StatusOK, claims)
	})

	e.Logger.Fatal(e.Start(":" + httpPort))
}
