package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/control"
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

func Init(ctrlChan chan control.Command, ps nats.PubSub, db *sql.DB) {
	// DB does not have any async loops, so does not need to signal readiness, it's ready as soon
	// as it's loaded, and has no internal components that would need to be stopped.
	// But it can fail while loading and that needs to be signalled.

	if err := migrateDB(db); err != nil {
		signal(ctrlChan, control.FAILED, fmt.Errorf("failed to migrate the database schema: %w", err))
		return
	}

	if err := registerStateRecorders(ps, db); err != nil {
		signal(ctrlChan, control.FAILED, err)
		return
	}
}

func migrateDB(db *sql.DB) error {
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

	dbLogger.Info("database schema migrated")
	return nil
}

// registerStateRecorders registers all async handlers needed for saving persistent state of the system.
func registerStateRecorders(ps nats.PubSub, db *sql.DB) error {
	if err := recorders.RegisterPlayerStateHandlers(ps, dbLogger, db); err != nil {
		return fmt.Errorf("failed to register player state handlers: %w", err)
	}

	dbLogger.Info("state persistence recorders registered")
	return nil
}

func signal(ctrlChan chan control.Command, state control.ComponentState, err error) {
	ctrlChan <- control.Command{
		Signal:    control.COMPONENT,
		Component: control.DB,
		State:     state,
		Err:       err,
	}
}
