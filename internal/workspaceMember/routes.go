package workspacemember

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
	sendInviteHandler := idempotencyMiddleware.RequireIdempotencyKey(
		idempotencyMiddleware.RequireIdempotency("workspace-invitations.create")(
			http.HandlerFunc(handler.SendInvite),
		),
	)

	mux.Handle(
		"POST /workspaces/invitations",
		jwtMiddleware.RequireAuth(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				principal, ok := middleware.GetPrincipal(r.Context())
				if !ok {
					response.Write(w, http.StatusUnauthorized, "authentication failed in routes")
					return
				}

				ctx := idempotency.WithIdempotencyScope(
					r.Context(),
					"workspace-invitations:create:"+principal.UserID,
				)

				sendInviteHandler.ServeHTTP(w, r.WithContext(ctx))
			}),
		),
	)
}
