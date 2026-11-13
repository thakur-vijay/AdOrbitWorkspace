package database

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/thakur-vijay/adorbit-go/apperrors"
)

func Connect(ctx context.Context, dbURL string) (*sql.DB, error) {
	if dbURL == "" {
		return nil, apperrors.NewConfigError("Database url is not set.", nil)
	}

	db, err := sql.Open("pgx", dbURL)

	if err != nil {
		return nil, apperrors.NewServiceUnavailable("Database connection failed", err)
	}

	if ctx != nil {
		if err := db.PingContext(ctx); err != nil {
			_ = db.Close()
			return nil, apperrors.NewServiceUnavailable("Database ping failed", err)
		}
	}

	return db, nil
}
