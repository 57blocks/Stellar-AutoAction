package db

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/57blocks/auto-action/server/internal/config"
	migs "github.com/57blocks/auto-action/server/internal/pkg/db/migration"
	pkgLog "github.com/57blocks/auto-action/server/internal/pkg/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var Conn *sql.DB

func Setup() error {
	if err := connect(); err != nil {
		return err
	}

	return migrateDB(Conn)
}

func connect() error {
	rds := config.Global.RDS

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=10",
		rds.Host, rds.Port, rds.User, rds.Password, rds.Database, rds.SSLMode,
	)

	var err error
	Conn, err = sql.Open("postgres", dsn)
	if err != nil {
		pkgLog.Logger.ERROR(fmt.Sprintf("setup database error: %s\n", err.Error()))
		return errors.New(fmt.Sprintf("setup database error: %s\n", err.Error()))
	}

	err = Conn.Ping()
	if err != nil {
		pkgLog.Logger.ERROR(fmt.Sprintf("connecting to database error: %s\n", err.Error()))
		return errors.New(fmt.Sprintf("connecting to database error: %s\n", err.Error()))
	}

	pkgLog.Logger.DEBUG("connected to database")

	return nil
}

func migrateDB(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		errMsg := fmt.Sprintf("new driver instance error: %s\n", err.Error())
		pkgLog.Logger.ERROR(errMsg)
		return errors.New(errMsg)
	}

	source, err := httpfs.New(http.FS(migs.Migrations), ".")
	if err != nil {
		errMsg := fmt.Sprintf("new embed source error: %s\n", err.Error())
		pkgLog.Logger.ERROR(errMsg)
		return errors.New(errMsg)
	}

	mig, err := migrate.NewWithInstance(
		"httpfs",
		source,
		"st3llar",
		driver,
	)
	if err != nil {
		errMsg := fmt.Sprintf("new migration instance error: %s\n", err.Error())
		pkgLog.Logger.ERROR(errMsg)
		return errors.New(errMsg)
	}

	var migErr error
	migErr = mig.Up()
	if migErr == nil || errors.Is(migErr, migrate.ErrNoChange) {
		return nil
	}

	dirtyErr := migrate.ErrDirty{}
	if errors.As(migErr, &dirtyErr) {
		lastSuccess := dirtyErr.Version - 1
		if err := mig.Force(lastSuccess); err != nil {
			errMsg := fmt.Sprintf(
				"force dirty version failed: %s\n",
				err.Error(),
			)
			pkgLog.Logger.ERROR(errMsg)

			return errors.New(errMsg)
		}
		if err := mig.Up(); err != nil {
			errMsg := fmt.Sprintf(
				"re-migrate dirty version failed: %s\n",
				err.Error(),
			)
			pkgLog.Logger.ERROR(errMsg)

			return errors.New(errMsg)
		}
		pkgLog.Logger.DEBUG(fmt.Sprintf(
			"re-migrate dirty version: %v successfully\n",
			dirtyErr.Version,
		))

		return nil
	}

	errMsg := fmt.Sprintf("migrate error: %s\n", migErr.Error())
	pkgLog.Logger.ERROR(errMsg)

	return errors.New(errMsg)
}
