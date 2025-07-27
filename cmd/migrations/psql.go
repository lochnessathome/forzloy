package migrations

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const psqlMigrationsPath = "file://migrations/psql"

func PsqlUp() error {
	dbUrl := os.Getenv("DATABASE_URL")
        if dbUrl == "" {
                return fmt.Errorf("Empty DATABASE_URL environment variable")
        }

	m, err := migrate.New(psqlMigrationsPath, dbUrl)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}
