package migrate

import (
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"

	"github.com/0div/fog/internal/db/postgres"
	"github.com/0div/fog/internal/store"
)

type MigrateOpts struct {
	TargetVersion int
}

func Migrate(storeName string, operation string, opts MigrateOpts) {
	slog.Info("Starting migration")

	var migrater store.Migrator

	switch storeName {
	case "postgres":
		pg := postgres.NewPostgresDB()
		pgMigrater, err := pg.GetMigrater()
		if err != nil {
			slog.Error("Failed to get migrater", "err", err)
			os.Exit(1)
		}
		migrater = pgMigrater
		defer migrater.Close()
	default:
		slog.Error("Unknown store, can't migrate")
		os.Exit(1)
	}

	var err error
	switch operation {
	case "up":
		err = migrater.Up()
	case "down":
		err = migrater.Down()
	case "list":
		var migrations []string
		migrations, err = migrater.List()
		if err != nil {
			break
		}
		slog.Info("", "migrations", migrations)
	case "version":
		var version uint
		var dirty bool
		version, dirty, err = migrater.Version()
		if err != nil {
			break
		}
		slog.Info("", "version", version, "dirty", dirty)

	case "force":
		err = migrater.Force(opts.TargetVersion)
		if err != nil {
			break
		}

	case "to":
		if opts.TargetVersion < 0 {
			os.Exit(1)
		}
		err = migrater.To(uint(opts.TargetVersion))
		if err != nil {
			break
		}
	}

	if err == migrate.ErrNoChange {
		slog.Warn("Already at the correct version, migration was skipped")
	} else if err == migrate.ErrNilVersion {
		slog.Warn("Migration is at nil version (no migrations have been performed)")
	} else if err != nil {
		slog.Error("Migration operation failed", "err", err)
		os.Exit(1)
	}

	slog.Info("Migration ended")
}
