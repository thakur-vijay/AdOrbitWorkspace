package workspace

import (
	"net/http"

	"github.com/thakur-vijay/adorbit-go/idempotency"
	"github.com/thakur-vijay/adorbit-go/middleware"
	"github.com/thakur-vijay/adorbit-go/response"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	jwtMiddleware *middleware.JWTMiddleware,
	idempotencyMiddleware *idempotency.Middleware,
) {
	createHandler := idempotencyMiddleware.RequireIdempotencyKey(
		idempotencyMiddleware.RequireIdempotency("workspace.create")(
			http.HandlerFunc(handler.Create),
		),
	)

	mux.Handle(
		"POST /workspace",
		jwtMiddleware.RequireAuth(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				principal, ok := middleware.GetPrincipal(r.Context())
				if !ok {
					response.Write(w, http.StatusUnauthorized, "authentication failed in routes")
					return
				}

				ctx := idempotency.WithIdempotencyScope(
					r.Context(),
					"workspace:create:"+principal.UserID,
				)

				createHandler.ServeHTTP(w, r.WithContext(ctx))
			}),
		),
	)
}
