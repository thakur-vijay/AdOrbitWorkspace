package workspace

import (
	"context"

	"github.com/gosimple/slug"
	"github.com/thakur-vijay/adorbit-go/apperrors"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, userID string, dto CreateWorkspaceDTO) (*Workspace, error) {
	workspaceSlug := slug.Make(dto.Name)
	workspace, err := s.repository.Create(ctx, dto, userID, workspaceSlug)

	if err != nil {
		return nil, apperrors.Classify(err)
	}

	return workspace, nil
}
