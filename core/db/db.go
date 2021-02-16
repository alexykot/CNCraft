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

func NewMigrator(db *sql.DB) (*migrate.Migrate, error) {
	s := bindata.Resource(AssetNames(), func(name string) (bytes []byte, e error) {
		return Asset(name)
	})

	if source, err := bindata.WithInstance(s); err != nil {
		return nil, err
	} else if driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "cncraft_schema_migrations",
	}); err != nil {
		return nil, err
	} else {
		return migrate.NewWithInstance("go-bindata", source, "auth", driver)
	}
}
