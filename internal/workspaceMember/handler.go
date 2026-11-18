package workspacemember

import (
	"context"
	"net/http"

	"github.com/thakur-vijay/adorbit-go/middleware"
	"github.com/thakur-vijay/adorbit-go/response"
)

type Service interface {
	SendInvite(ctx context.Context, dto SendInviteDTO, userID string) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Send Invite godoc
// @Summary Send workspace invite to email
// @Description Sends an invitation link on email.
// @Description Requires an Idempotency-Key to safely handle retries.
// @Tags Workspaces
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param Idempotency-Key header string true "Unique UUID for request idempotency"
// @Param        request body SendInviteDTO true "Workspace details"
// @Success      201 {object} response.Response "Invitation sent successfully"
// @Failure      400 {object} response.Response "Invalid request payload"
// @Failure      401 {object} response.Response "Authentication failed"
// @Failure      422 {object} response.Response "Validation failed"
// @Failure      429 {object} response.Response "Rate limit exceeded"
// @Failure      500 {object} response.Response "Internal server error"
// @Failure      503 {object} response.Response "Service unavailable"
// @Router       /workspaces/invitations [post]
func (h *Handler) SendInvite(w http.ResponseWriter, r *http.Request) {
	var dto SendInviteDTO
	principal, ok := middleware.GetPrincipal(r.Context())

	if !ok {
		response.Write(w, http.StatusUnauthorized, "authentication failed in handler")
		return
	}

	ctx := r.Context()

	if err := response.DecodeJSONAndValidate(r, &dto); err != nil {
		response.WriteError(w, err)
		return
	}

	err := h.service.SendInvite(ctx, dto, principal.UserID)

	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.Write(w, http.StatusCreated, "Invitation sent successfully")
}
