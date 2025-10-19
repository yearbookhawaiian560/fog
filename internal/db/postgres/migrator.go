package postgres

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type PostgresDBMigrater struct {
	m     *migrate.Migrate
	close func() error
}

func (mig *PostgresDBMigrater) Up() error {
	return mig.m.Up()
}

func (mig *PostgresDBMigrater) Down() error {
	return mig.m.Down()
}

func (mig *PostgresDBMigrater) Version() (uint, bool, error) {
	return mig.m.Version()
}

func (mig *PostgresDBMigrater) To(version uint) error {
	return mig.m.Migrate(version)
}

func (mig *PostgresDBMigrater) Force(version int) error {
	return mig.m.Force(version)
}

func (mig *PostgresDBMigrater) List() ([]string, error) {
	dirEntries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	migrationFiles := make([]string, 0)
	for _, entry := range dirEntries {
		migrationFiles = append(migrationFiles, entry.Name())
	}
	return migrationFiles, nil
}

func (mig *PostgresDBMigrater) Close() error {
	return mig.close()
}

func (mig *PostgresDBMigrater) SetLogger(logger migrate.Logger) {
	mig.m.Log = logger
}

func (pgs *PostgresDB) GetMigrater() (*PostgresDBMigrater, error) {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to open Postgres migrations iofs: %w", err)
	}

	connString := pgs.DB.Config().ConnConfig.ConnString()
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres db with pgx driver: %w", err)
	}
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		//closeError := driver.Close()
		return nil, fmt.Errorf("failed to open postgres migration: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs", d,
		pgs.dbname, driver)

	if err != nil {
		return nil, fmt.Errorf("failed to create Postgres migrate instance: %w", err)
	}

	close := func() error {
		err1, err2 := m.Close()
		if err1 != nil || err2 != nil {
			return fmt.Errorf("source close error: %v, driver close error: %v", err1, err2)
		}
		return nil
	}

	return &PostgresDBMigrater{
		m:     m,
		close: close,
	}, nil
}
