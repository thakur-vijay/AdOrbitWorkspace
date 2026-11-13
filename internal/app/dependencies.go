package app

import (
	"database/sql"
	"net/http"

	"github.com/adorbit/workspace/internal/api"
	"github.com/adorbit/workspace/internal/config"
	"github.com/adorbit/workspace/internal/workspace"
	"github.com/thakur-vijay/adorbit-go/idempotency"
	"github.com/thakur-vijay/adorbit-go/middleware"
)

func BuildRouter(db *sql.DB, cfg *config.Config) (http.Handler, error) {
	jwtMiddleware, err := middleware.NewJWTMiddleware(cfg.JWTPublicKey, "adorbit-workspace")

	if err != nil {
		return nil, err
	}

	idempotencyRepository := idempotency.NewPostgresRepository(db)

	idempotencyService, err := idempotency.NewService(
		idempotencyRepository,
		cfg.IdempotencyConfig(),
	)

	if err != nil {
		return nil, err
	}

	idempotencyMiddleware := idempotency.NewHTTPMiddleware(idempotencyService)
	workspaceRepository := workspace.NewRepository(db)
	workspaceService := workspace.NewService(workspaceRepository)
	workspaceHandler := workspace.NewHandler(workspaceService)
	router := api.NewRouter(api.Dependencies{WorkspaceHandler: workspaceHandler, JWTMiddleware: jwtMiddleware, IdempotencyMiddleware: idempotencyMiddleware})
	return router, nil
}
