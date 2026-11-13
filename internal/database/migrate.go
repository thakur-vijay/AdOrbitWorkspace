package database

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/thakur-vijay/adorbit-go/apperrors"
	"github.com/thakur-vijay/adorbit-go/idempotency"
)

func Migrate(dbURL string) error {
	if dbURL == "" {
		return apperrors.NewConfigError("Database URL is not set.", nil)
	}

	workspaceMigration, err := migrate.New("file://internal/database/migrations", dbURL)

	if err != nil {
		return apperrors.NewServiceUnavailable("workspace migration setup failed", err)
	}

	err = workspaceMigration.Up()
	_, closeErr := workspaceMigration.Close()

	if err != nil && err != migrate.ErrNoChange {
		return apperrors.NewServiceUnavailable("workspace migration failed", err)
	}

	if closeErr != nil {
		return apperrors.NewServiceUnavailable("workspace close migration failed", err)
	}

	source, err := iofs.New(idempotency.Migrations, "migrations")

	if err != nil {
		return apperrors.NewServiceUnavailable("idempotency migration source setup failed", err)
	}

	migrationDB, err := sql.Open("pgx", dbURL)

	if err != nil {
		source.Close()
		return apperrors.NewServiceUnavailable("idempotency database connection failed.", err)
	}

	if err := migrationDB.Ping(); err != nil {
		migrationDB.Close()
		source.Close()
		return apperrors.NewServiceUnavailable("idempotency database ping failed.", err)
	}

	dbDriver, err := migratepostgres.WithInstance(migrationDB, &migratepostgres.Config{
		MigrationsTable: "idempotency_schema_migrations",
	})

	if err != nil {
		migrationDB.Close()
		source.Close()
		return apperrors.NewServiceUnavailable("idempotency database migration setup failed", err)
	}

	idempotencyMigration, err := migrate.NewWithInstance(
		"iofs",
		source,
		"adorbit_workspace",
		dbDriver,
	)

	if err != nil {
		source.Close()
		dbDriver.Close()
		return apperrors.NewServiceUnavailable(
			"idempotency migration setup failed",
			err,
		)
	}

	defer idempotencyMigration.Close()

	err = idempotencyMigration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return apperrors.NewServiceUnavailable(
			"idempotency migration failed",
			err,
		)
	}

	return nil
}
