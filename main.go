package main

import (
	"os"

	"billing/cmd/migrations"
	"billing/internal/handlers"
	"billing/internal/psql"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const httpPort = "8080"

type Validator struct {
	validator *validator.Validate
}

func main() {
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	err := migrations.PsqlUp()
	if err != nil {
		e.Logger.Error(err)
		os.Exit(1)
	}

	pgPool, err := psql.New()
	if err != nil {
		e.Logger.Error(err)
		os.Exit(1)
	}

	h := handlers.New(pgPool)

	/*	e.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
		})

		e.GET("/time", func(c echo.Context) error {
			var tm string
			err = pgPool.QueryRow(ctx, "SELECT NOW()::text;").Scan(&tm)
			if err != nil {
				e.Logger.Error(err)
			}

			return c.JSON(http.StatusOK, struct{ Time string }{Time: tm})
		}) */

	gAuth := e.Group("/api/auth")

	gAuth.GET("/create-token", handlers.AuthCreateToken)
	gAuth.GET("/validate-token/:ts", handlers.AuthValidateToken)

	gAuth.POST("/register", h.AuthRegister)

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func (cv *Validator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if err != nil {
		return err
	}
	return nil
}
