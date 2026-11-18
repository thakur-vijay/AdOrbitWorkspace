package workspacemember

import (
	"context"
	"time"

	"github.com/thakur-vijay/adorbit-go/apperrors"
)

type Repository interface {
	// Workspace Members
	AddMember(ctx context.Context, workspaceID, userID, role string) error
	FetchMembers(ctx context.Context, workspaceID string) ([]WorkspaceMember, error)
	FetchMember(ctx context.Context, workspaceID, userID string) (*WorkspaceMember, error)
	UpdateMemberRole(ctx context.Context, workspaceID, userID, role string) error
	RemoveMember(ctx context.Context, workspaceID, userID string) error
	FetchUserWorkspaces(ctx context.Context, userID string) ([]WorkspaceMember, error)

	// Workspace Invitations
	CreateInvitation(ctx context.Context, invitation WorkspaceInvitation) error
	FetchInvitationByID(ctx context.Context, invitationID string) (*WorkspaceInvitation, error)
	FetchInvitationByTokenHash(ctx context.Context, tokenHash string) (*WorkspaceInvitation, error)
	FetchPendingInvitations(ctx context.Context, workspaceID string) ([]WorkspaceInvitation, error)
	RevokeInvitation(ctx context.Context, invitationID string) error
	AcceptInvitation(ctx context.Context, invitationID string, userID string) error
	ExpireInvitations(ctx context.Context) error
}

type service struct {
	repository Repository
}

var _ Service = (*service)(nil)

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) SendInvite(ctx context.Context, dto SendInviteDTO, userID string) error {
	//generate token and its hash, create link from token and save hash in DTO
	tokenHash := ""
	workspaceInvitationDTO := WorkspaceInvitation{
		WorkspaceID: dto.WorkspaceID,
		Email:       dto.Email,
		TokenHash:   tokenHash,
		InvitedBy:   userID,
		ExpiresAt:   time.Now().AddDate(0, 0, 7),
	}

	err := s.repository.CreateInvitation(ctx, workspaceInvitationDTO)

	// send email

	if err != nil {
		return apperrors.Classify(err)
	}

	return nil
}
