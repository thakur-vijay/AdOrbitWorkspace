package workspace

import (
	"context"
	"database/sql"
	"fmt"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

var _ Repository = (*repository)(nil)

func (r *repository) Create(ctx context.Context, dto CreateWorkspaceDTO, userID, slug string) (*Workspace, error) {
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

func (r *repository) Fetch(ctx context.Context, userID string) ([]Workspace, error) {
	query := `
        SELECT
            id,
            name,
            slug,
            type,
            country,
            timezone,
            currency,
            industry,
            website,
            description,
            status,
            created_by,
            created_at,
            updated_at
        FROM workspaces
        WHERE created_by = $1
          AND status = 'active'
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("fetch workspaces: %w", err)
	}
	defer rows.Close()

	workspaces := make([]Workspace, 0)

	for rows.Next() {
		var workspace Workspace

		err := rows.Scan(
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
		if err != nil {
			return nil, fmt.Errorf("scan workspace: %w", err)
		}

		workspaces = append(workspaces, workspace)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate workspaces: %w", err)
	}

	return workspaces, nil
}
