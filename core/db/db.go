package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/db/recorders"
	"github.com/alexykot/cncraft/core/nats"
)

var dbLogger *zap.Logger

func New(url string, isDebug bool, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("could not initialize database: %w", err)
	}

	dbLogger = logger
	boil.DebugMode = isDebug

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

	if err = migrator.Up(); err == migrate.ErrNoChange {
		ver, _, _ := migrator.Version()
		dbLogger.Info(fmt.Sprintf("database schema up to date at version %d", ver))
	} else if err != nil {
		return fmt.Errorf("failed to migrate DB: %w", err)
	} else {
		ver, _, _ := migrator.Version()
		dbLogger.Info(fmt.Sprintf("database schema migrated to latest version %d", ver))
	}

	if err = source.Close(); err != nil {
		return fmt.Errorf("failed to close migrator source: %w", err)
	}

	return nil
}

// RegisterStateRecorders registers all async handlers needed for saving persistent state of the system.
func RegisterStateRecorders(ps nats.PubSub, db *sql.DB) error {
	if err := recorders.RegisterPlayerStateHandlers(ps, dbLogger, db); err != nil {
		return fmt.Errorf("failed to register player state handlers: %w", err)
	}

	dbLogger.Info("state persistence recorders registered")

	return nil
}
