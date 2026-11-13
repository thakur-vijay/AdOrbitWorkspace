package workspace

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, dto CreateWorkspaceDTO, userID, slug string) (*Workspace, error) {
	var workspace Workspace
	err := r.db.QueryRowContext(ctx, `INSERT INTO workspaces (name, slug, type, country, timezone, currency, industry, website, description, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, name, slug, type, country, timezone, currency, industry, website, description, status, created_by, created_at, updated_at`, dto.Name, slug, dto.Type, dto.Country, dto.Timezone, dto.Currency, dto.Industry, dto.Website, dto.Description, userID).Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.Slug,
		&workspace.Type,
		&workspace.Country,
		&workspace.Timezone,
		&workspace.Currency,
		&workspace.Industry,
		&workspace.Website,
		&workspace.Description,
		&workspace.Status,
		&workspace.CreatedBy,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)

	return &workspace, err
}
