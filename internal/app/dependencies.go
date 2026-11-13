package app

import (
	"database/sql"
	"net/http"

	"github.com/adorbit/workspace/internal/api"
	"github.com/thakur-vijay/adorbit-go/middleware"
)

func BuildRouter(db *sql.DB, publicKeyPEM string) (http.Handler, error) {
	jwtMiddleware, err := middleware.NewJWTMiddleware(publicKeyPEM, "adorbit-workspace")

	if err != nil {
		return nil, err
	}

	router := api.NewRouter(api.Dependencies{JWTMiddleware: jwtMiddleware})
	return router, nil
}
