package workspacemember

import "time"

type WorkspaceMember struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	UserID      string    `json:"user_id"`
	Role        string    `json:"role"`
	JoinedAt    time.Time `json:"joined_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WorkspaceInvitation struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	TokenHash   string    `json:"token_hash"`
	Status      string    `json:"status"`
	InvitedBy   string    `json:"invited_by"`
	ExpiresAt   time.Time `json:"expires_at"`
	AcceptedAt  time.Time `json:"accepted_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SendInviteDTO struct {
	Email string `json:"email" validate:"required,email"`
}

// Invitation token should be cryptographically random.

// Store token hash, not the raw token.

// Invitation should expire (e.g. 7 days).

// Token must be single-use.

// Verify authenticated user's email before accepting.

// Never allow users to assign themselves a higher role.
