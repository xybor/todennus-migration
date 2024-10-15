package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/xybor/x/xcontext"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Up(ctx context.Context, db *sql.DB, migrationPath string) error {
	m, err := getMigrator(db, migrationPath)
	if err != nil {
		return err
	}

	err = m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	xcontext.Logger(ctx).Info("migrate up postgres successfully")
	return nil
}

func Down(ctx context.Context, db *sql.DB, migrationPath string, n int) error {
	m, err := getMigrator(db, migrationPath)
	if err != nil {
		return err
	}

	err = m.Steps(-n)
	if !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	xcontext.Logger(ctx).Info("migrate down postgres successfully")
	return nil
}

func getMigrator(db *sql.DB, migrationPath string) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	return migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationPath),
		"postgres",
		driver,
	)
}
