package main

import (
	"context"
	"net/http"
	"os"

	"billing/internal/psql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const httpPort = "80"

func main() {
	ctx := context.Background()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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


	e.Logger.Fatal(e.Start(":" + httpPort))
}
