package postgres

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxvec "github.com/pgvector/pgvector-go/pgx"

	"github.com/0div/fog/internal/cfg"
	apierror "github.com/0div/fog/internal/errors"
)

type PostgresDB struct {
	dbname string
	DB     *pgxpool.Pool
	Q      *Queries
}

func (pdb *PostgresDB) IsPgErr(err error) (bool, *pgconn.PgError) {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return true, pgErr
		}
	}
	return false, &pgconn.PgError{}
}

func (pdb *PostgresDB) Pg2ApiErr(err error) (bool, *apierror.APIError) {
	isPgErr, pgErr := pdb.IsPgErr(err)
	if isPgErr {
		// https://github.com/jackc/pgerrcode/blob/masetr/errcode.go
		if pgerrcode.IsDataException(pgErr.Code) || pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return true, &apierror.APIError{
				ErrorCode: http.StatusBadRequest,
				Detail:    pgErr.Message,
				RawError:  pgErr,
			}
		}
		return true, &apierror.APIError{
			ErrorCode: http.StatusInternalServerError,
			RawError:  pgErr,
			Detail:    "database error",
		}
	}

	if err == pgx.ErrNoRows {
		return true, &apierror.APIError{
			ErrorCode: http.StatusBadRequest,
			Detail:    "resource not found",
			RawError:  pgErr,
		}
	}

	return false, &apierror.APIError{
		ErrorCode: http.StatusInternalServerError,
		RawError:  err,
	}
}

func BuildConnectionDSN(dbname string) string {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable",
		cfg.Str("POSTGRES_HOST"), cfg.Str("POSTGRES_PORT"), cfg.Str("POSTGRES_USER"), dbname)

	return dsn
}

func NewPostgresDB() *PostgresDB {
	dbname := cfg.Str("POSTGRES_DB")
	dsn := BuildConnectionDSN(dbname)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		slog.Error("error parsing Postgres DSN", "err", err)
		panic(err)
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return pgxvec.RegisterTypes(ctx, conn)
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		slog.Error("error creating a Postgres connection pool", "err", err)
		panic(err)
	}

	return &PostgresDB{
		dbname: dbname,
		DB:     dbpool,
		Q:      New(dbpool),
	}
}

func NewPostgresTestStore() (store *PostgresDB, cleanup func()) {
	ctx, cancel := context.WithCancel(context.Background())
	dbnameOriginal := cfg.Str("POSTGRES_DB")
	// We first need to authenticate against the ordinary dbname
	// as we can't connect without specifying a database. Then from within that
	// database we can create a temporary test database, and connect to that.
	dbBootstrap, err := pgxpool.New(ctx, BuildConnectionDSN(dbnameOriginal))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to postgres bootstrap store", "err", err)
		panic(err)
	}

	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	dbname := dbnameOriginal + "_test_" + strings.ToLower(fmt.Sprintf("%X", b))
	_, err = dbBootstrap.Exec(ctx, `CREATE DATABASE `+dbname+`;`)
	if err != nil {
		slog.ErrorContext(ctx, "Cloud not create postgres test DB", "err", err)
		panic(err)
	}

	db, err := pgxpool.New(ctx, BuildConnectionDSN(dbname))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to postgrestest store", "err", err)
		panic(err)
	}

	cleanup = func() {
		db.Close()

		_, err = dbBootstrap.Exec(ctx, fmt.Sprintf("DROP DATABASE %s WITH (FORCE);", dbname))
		if err != nil {
			slog.Error("Couldn't DROP db", "err", err)
			panic(err)
		}
		dbBootstrap.Close()
		cancel()

	}

	store = &PostgresDB{
		dbname: dbname,
		DB:     db,
		Q:      New(db),
	}

	// TODO
	//mig, err := store.GetMigrater()
	//if err != nil {
	//	slog.Log(ctx, loglevel.Fatal, "Failed to create migrater for postgres store", logattr.Error(err))
	//	cleanup()
	//	os.Exit(1)
	//}
	//defer mig.Close()

	//err = mig.Up()
	//if err != nil {
	//	slog.Log(ctx, loglevel.Fatal, "Failed to migrate Postgres test store", logattr.Error(err))
	//	cleanup()
	//	os.Exit(1)
	//}

	return store, cleanup
}
