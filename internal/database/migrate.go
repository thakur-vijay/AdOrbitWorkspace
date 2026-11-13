package database

import (
	"github.com/golang-migrate/migrate"
	"github.com/thakur-vijay/adorbit-go/apperrors"
)

func Migrate(dbURL string) error {
	if dbURL == "" {
		return apperrors.NewConfigError("Database URL is not set.", nil)
	}

	m, err := migrate.New("file://internal/database/migrations", dbURL)

	if err != nil {
		return apperrors.NewServiceUnavailable("database migration setup failed", err)
	}

	defer m.Close()

	err = m.Up()

	if err != nil && err != migrate.ErrNoChange {
		return apperrors.NewServiceUnavailable("database migration failed.", err)
	}

	return nil
}
