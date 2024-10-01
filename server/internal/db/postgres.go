package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/57blocks/auto-action/server/internal/config"
	migs "github.com/57blocks/auto-action/server/internal/db/migration"
	"github.com/57blocks/auto-action/server/internal/pkg/errorx"
	"github.com/57blocks/auto-action/server/internal/third-party/logx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/lib/pq"
	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Inst *Instance

type Instance struct {
	*gorm.DB
}

func (c *Instance) Conn(ctx context.Context) *gorm.DB {
	return c.DB.WithContext(ctx)
}

func Setup() error {
	if err := connect(); err != nil {
		return err
	}

	return migrateDB(Inst.DB)
}

func connect() error {
	rds := config.GlobalConfig.RDS

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=10",
		rds.Host, rds.Port, rds.User, rds.Password, rds.Database, rds.SSLMode,
	)

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("setup database error: %s", err.Error()))
	}
	sqlDB.SetMaxIdleConns(4) // by default: 2
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)

	// db: *gorm.DB
	// db.ConnPool: {gorm.ConnPool | *gorm.PreparedStmtDB}
	db, err := gorm.Open(
		pgDriver.Open(dsn),
		&gorm.Config{
			DisableAutomaticPing:   false,
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
			ConnPool: sqlDB,
		},
	)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("connecting to database error: %s", err.Error()))
	}

	Inst = new(Instance)
	Inst.DB = db

	return nil
}

func migrateDB(db *gorm.DB) error {
	instance, err := db.DB()
	if err != nil {
		return errorx.Internal(fmt.Sprintf("new DB instance error: %s", err.Error()))
	}

	driver, err := postgres.WithInstance(instance, &postgres.Config{
		MigrationsTable: "migration_version",
	})
	if err != nil {
		return errorx.Internal(fmt.Sprintf("new driver instance error: %s", err.Error()))
	}

	source, err := httpfs.New(http.FS(migs.Migrations), ".")
	if err != nil {
		return errorx.Internal(fmt.Sprintf("new embed source error: %s", err.Error()))
	}

	mig, err := migrate.NewWithInstance(
		"httpfs",
		source,
		"autoaction",
		driver,
	)
	if err != nil {
		return errorx.Internal(fmt.Sprintf("new migration instance error: %s", err.Error()))
	}

	migErr := mig.Up()
	if migErr == nil || errors.Is(migErr, migrate.ErrNoChange) {
		return nil
	}

	dirtyErr := migrate.ErrDirty{}
	if errors.As(migErr, &dirtyErr) {
		lastSuccess := dirtyErr.Version - 1
		if err := mig.Force(lastSuccess); err != nil {
			return errorx.Internal(fmt.Sprintf(
				"force dirty version failed: %s",
				err.Error(),
			))
		}
		if err := mig.Up(); err != nil {
			return errorx.Internal(fmt.Sprintf(
				"re-migrate dirty version failed: %s",
				err.Error(),
			))
		}
		logx.Logger.DEBUG(fmt.Sprintf(
			"re-migrate dirty version: %v successfully",
			dirtyErr.Version,
		))

		return nil
	}

	return errorx.Internal(fmt.Sprintf("migrate error: %s", migErr.Error()))
}
