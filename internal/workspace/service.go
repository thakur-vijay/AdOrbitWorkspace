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

	FetchByID(
		ctx context.Context,
		workspaceID string,
	) (*Workspace, error)

	Update(
		ctx context.Context,
		dto UpdateWorkspaceDTO,
		workspaceID string,
		slug *string,
	) (*Workspace, error)

	Delete(
		ctx context.Context,
		workspaceID string,
	) error
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

func (s *service) FetchByID(
	ctx context.Context,
	workspaceID string,
) (*Workspace, error) {
	workspace, err := s.repository.FetchByID(ctx, workspaceID)
	if err != nil {
		return nil, apperrors.Classify(err)
	}

	return workspace, nil
}

func (s *service) Update(
	ctx context.Context,
	dto UpdateWorkspaceDTO,
	workspaceID string,
) (*Workspace, error) {
	if dto.Name == nil && dto.Country == nil && dto.Timezone == nil && dto.Currency == nil && dto.Industry == nil && dto.Website == nil && dto.Description == nil {
		return nil, apperrors.NewBadRequest("at least one workspace field must be provided", nil)
	}

	var workspaceSlug *string

	if dto.Name != nil {
		generatedSlug := slug.Make(*dto.Name)
		workspaceSlug = &generatedSlug
	}

	workspace, err := s.repository.Update(ctx, dto, workspaceID, workspaceSlug)
	if err != nil {
		return nil, apperrors.Classify(err)
	}

	return workspace, nil
}

func (s *service) Delete(
	ctx context.Context,
	workspaceID string,
) error {

	err := s.repository.Delete(ctx, workspaceID)

	return apperrors.Classify(err)
}
