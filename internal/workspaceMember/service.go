package workspacemember

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"strings"
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
	repository  Repository
	mailService *InvitationMailService
	appDomain   string
}

var _ Service = (*service)(nil)

func NewService(repository Repository, mailService *InvitationMailService, appDomain string) Service {
	return &service{
		repository:  repository,
		mailService: mailService,
		appDomain:   appDomain,
	}
}

func (s *service) SendInvite(ctx context.Context, dto SendInviteDTO, userID string) error {
	//generate token and its hash, create link from token and save hash in DTO
	token, err := generateInvitationToken()
	if err != nil {
		return apperrors.Classify(err)
	}

	tokenHash := hashInvitationToken(token)

	expiresAt := time.Now().AddDate(0, 0, 7)
	workspaceInvitationDTO := WorkspaceInvitation{
		WorkspaceID: dto.WorkspaceID,
		Email:       dto.Email,
		TokenHash:   tokenHash,
		InvitedBy:   userID,
		ExpiresAt:   expiresAt,
	}

	err = s.repository.CreateInvitation(ctx, workspaceInvitationDTO)

	if err != nil {
		return apperrors.Classify(err)
	}

	invitationLink := fmt.Sprintf(
		"%s/invitations/accept?token=%s",
		strings.TrimRight(s.appDomain, "/"),
		url.QueryEscape(token),
	)

	mailDTO := InvitationEmailDTO{
		Email:          dto.Email,
		WorkspaceName:  dto.WorkspaceName,
		InviterName:    dto.InviterName,
		Role:           Member,
		InvitationLink: invitationLink,
		ExpiresAt:      expiresAt,
	}
	err = s.mailService.Send(ctx, mailDTO)

	return err
}

func generateInvitationToken() (string, error) {
	token := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, token); err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil
}

func hashInvitationToken(token string) string {
	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}
