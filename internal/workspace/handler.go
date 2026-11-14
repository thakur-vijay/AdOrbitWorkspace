package workspace

import (
	"net/http"

	"github.com/thakur-vijay/adorbit-go/apperrors"
	"github.com/thakur-vijay/adorbit-go/middleware"
	"github.com/thakur-vijay/adorbit-go/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create godoc
// @Summary Create a new workspace
// @Description Creates a new advertiser or publisher workspace for the authenticated user.
// @Description Requires an Idempotency-Key to safely handle retries.
// @Tags Workspaces
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Param Idempotency-Key header string true "Unique UUID for request idempotency"
// @Param        request body CreateWorkspaceDTO true "Workspace details"
// @Success      201 {object} response.Response "Workspace created successfully"
// @Failure      400 {object} response.Response "Invalid request payload"
// @Failure      401 {object} response.Response "Authentication failed"
// @Failure      409 {object} response.Response "Workspace conflict"
// @Failure      422 {object} response.Response "Validation failed"
// @Failure      429 {object} response.Response "Rate limit exceeded"
// @Failure      500 {object} response.Response "Internal server error"
// @Failure      503 {object} response.Response "Service unavailable"
// @Router       /workspace [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var dto CreateWorkspaceDTO
	principal, ok := middleware.GetPrincipal(r.Context())

	if !ok {
		response.Write(w, http.StatusUnauthorized, "authentication failed in handler")
		return
	}

	ctx := r.Context()

	if err := response.DecodeJSONAndValidate(r, &dto); err != nil {
		response.WriteError(w, apperrors.Classify(err))
		return
	}

	workspace, err := h.service.Create(ctx, principal.UserID, dto)

	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteItem(w, http.StatusCreated, "Workspace created successfully", workspace)
}
