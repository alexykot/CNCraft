package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

func New(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("could not initialize database: %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	s := bindata.Resource(AssetNames(), func(name string) (bytes []byte, e error) {
		return Asset(name)
	})

	source, err := bindata.WithInstance(s)
	if err != nil {
		return fmt.Errorf("failed to create bindata source: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{MigrationsTable: "cncraft_schema_migrations"})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance("go-bindata", source, "auth", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err = migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to migrate DB: %w", err)
	}

	if err = source.Close(); err != nil {
		return fmt.Errorf("failed to close migrator source: %w", err)
	}

	return nil
}
