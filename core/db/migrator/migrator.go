package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewFileMigrator(db *sql.DB) (*migrate.Migrate, error) {
	var migrator *migrate.Migrate

	if driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "cncraft_schema_migrations",
	}); err != nil {
		return nil, fmt.Errorf("failed to instantiate DB driver: %w", err)
	} else if migrator, err = migrate.NewWithDatabaseInstance("file://schema/", "cncraft", driver); err != nil {
		return nil, fmt.Errorf("failed to instantiate migrator: %w", err)
	}

	return migrator, nil
}

func main() {
	var dbUrlFlag string

	flag.StringVar(&dbUrlFlag, "db-url", "", "DB URL")
	flag.Parse()

	var err error
	var db *sql.DB
	if db, err = sql.Open("postgres", dbUrlFlag); err != nil {
		log.Fatalf("failed to instantiate the DB: %v", err)
	}

	migrator, err := NewFileMigrator(db)
	if err != nil {
		log.Fatalf("failed to create migrator: %v", err)
	}
	defer migrator.Close() // ignore the error, it's end of script anyway

	if err = migrator.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to migrate upwards: %v", err)
	}
}
