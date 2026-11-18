package workspacemember

import (
	"context"
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}
func (r *repository) AddMember(ctx context.Context, workspaceID, userID, role string) error {
	return nil
}
func (r *repository) FetchMembers(ctx context.Context, workspaceID string) ([]WorkspaceMember, error) {
	return nil, nil
}
func (r *repository) FetchMember(ctx context.Context, workspaceID, userID string) (*WorkspaceMember, error) {
	return nil, nil
}
func (r *repository) UpdateMemberRole(ctx context.Context, workspaceID, userID, role string) error {
	return nil
}
func (r *repository) RemoveMember(ctx context.Context, workspaceID, userID string) error {
	return nil
}
func (r *repository) FetchUserWorkspaces(ctx context.Context, userID string) ([]WorkspaceMember, error) {
	return nil, nil
}

// Workspace Invitations
func (r *repository) CreateInvitation(
	ctx context.Context,
	invitation WorkspaceInvitation,
) error {
	query := `
        INSERT INTO workspace_invitations (
            workspace_id,
            email,
            role,
            token_hash,
            status,
            invited_by,
            expires_at
        )
        VALUES ($1, LOWER(TRIM($2)), $3, $4, 'pending', $5, $6)

        ON CONFLICT (workspace_id, (LOWER(email)))
        DO UPDATE SET
            email = EXCLUDED.email,
            role = EXCLUDED.role,
            token_hash = EXCLUDED.token_hash,
            status = 'pending',
            invited_by = EXCLUDED.invited_by,
            expires_at = EXCLUDED.expires_at,
            accepted_at = NULL,
            updated_at = NOW()
    `

	_, err := r.db.ExecContext(
		ctx,
		query,
		invitation.WorkspaceID,
		invitation.Email,
		invitation.Role,
		invitation.TokenHash,
		invitation.InvitedBy,
		invitation.ExpiresAt,
	)

	return err
}

func (r *repository) FetchInvitationByID(ctx context.Context, invitationID string) (*WorkspaceInvitation, error) {
	return nil, nil
}
func (r *repository) FetchInvitationByTokenHash(ctx context.Context, tokenHash string) (*WorkspaceInvitation, error) {
	return nil, nil
}
func (r *repository) FetchPendingInvitations(ctx context.Context, workspaceID string) ([]WorkspaceInvitation, error) {
	return nil, nil
}
func (r *repository) RevokeInvitation(ctx context.Context, invitationID string) error {
	return nil
}
func (r *repository) AcceptInvitation(ctx context.Context, invitationID string, userID string) error {
	return nil
}
func (r *repository) ExpireInvitations(ctx context.Context) error {
	return nil
}
