package workspace

import (
	"context"

	"github.com/gosimple/slug"
	"github.com/thakur-vijay/adorbit-go/apperrors"
)

type Repository interface {
	Create(
		ctx context.Context,
		dto CreateWorkspaceDTO,
		userID string,
		slug string,
	) (*Workspace, error)

	Fetch(
		ctx context.Context,
		userID string,
	) ([]Workspace, error)
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

func (s *service) Create(ctx context.Context, userID string, dto CreateWorkspaceDTO) (*Workspace, error) {
	workspaceSlug := slug.Make(dto.Name)
	workspace, err := s.repository.Create(ctx, dto, userID, workspaceSlug)

	if err != nil {
		return nil, apperrors.Classify(err)
	}

	return workspace, nil
}

func (s *service) Fetch(ctx context.Context, userID string) ([]Workspace, error) {
	workspaces, err := s.repository.Fetch(ctx, userID)

	if err != nil {
		return nil, apperrors.Classify(err)
	}

	return workspaces, nil
}
