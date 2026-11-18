package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

var _ Repository = (*repository)(nil)

func (r *repository) Create(
	ctx context.Context,
	dto CreateWorkspaceDTO,
	userID string,
	slug string,
) (*Workspace, error) {

	// 1. Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// 2. Rollback automatically if transaction is not committed
	defer tx.Rollback()

	// 3. Create workspace
	workspace, err := createWorkspaceTx(
		ctx,
		tx,
		dto,
		userID,
		slug,
	)
	if err != nil {
		return nil, err
	}

	// 4. Add creator as workspace owner
	err = addWorkspaceOwnerTx(
		ctx,
		tx,
		workspace.ID,
		userID,
	)
	if err != nil {
		return nil, err
	}

	// 5. Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return workspace, nil
}

func createWorkspaceTx(
	ctx context.Context,
	tx *sql.Tx,
	dto CreateWorkspaceDTO,
	userID string,
	slug string,
) (*Workspace, error) {

	var workspace Workspace

	query := `
        INSERT INTO workspaces (
            name,
            slug,
            type,
            country,
            timezone,
            currency,
            industry,
            website,
            description,
            created_by
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING
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
    `

	err := tx.QueryRowContext(
		ctx,
		query,
		dto.Name,
		slug,
		dto.Type,
		dto.Country,
		dto.Timezone,
		dto.Currency,
		dto.Industry,
		dto.Website,
		dto.Description,
		userID,
	).Scan(
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
		return nil, err
	}

	return &workspace, nil
}

func addWorkspaceOwnerTx(
	ctx context.Context,
	tx *sql.Tx,
	workspaceID string,
	userID string,
) error {

	query := `
        INSERT INTO workspace_members (
            workspace_id,
            user_id,
            role
        )
        VALUES ($1, $2, $3)
    `

	_, err := tx.ExecContext(
		ctx,
		query,
		workspaceID,
		userID,
		"owner",
	)

	return err
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

func (r *repository) FetchByID(ctx context.Context, workspaceID string) (*Workspace, error) {
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
        WHERE id = $1
    `

	var workspace Workspace
	err := r.db.QueryRowContext(ctx, query, workspaceID).Scan(
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

func (r *repository) Update(
	ctx context.Context,
	dto UpdateWorkspaceDTO,
	workspaceID string,
	slug *string,
) (*Workspace, error) {
	var (
		query  strings.Builder
		args   []any
		fields []string
	)

	if dto.Name == nil && dto.Country == nil && dto.Timezone == nil && dto.Currency == nil && dto.Industry == nil && dto.Website == nil && dto.Description == nil && slug == nil {
		return nil, errors.New("no workspace fields provided for update")
	}

	query.WriteString("UPDATE workspaces SET ")

	if dto.Name != nil {
		args = append(args, *dto.Name)
		fields = append(fields, fmt.Sprintf("name = $%d", len(args)))
	}

	if slug != nil {
		args = append(args, *slug)
		fields = append(fields, fmt.Sprintf("slug = $%d", len(args)))
	}

	if dto.Country != nil {
		args = append(args, *dto.Country)
		fields = append(fields, fmt.Sprintf("country = $%d", len(args)))
	}

	if dto.Timezone != nil {
		args = append(args, *dto.Timezone)
		fields = append(fields, fmt.Sprintf("timezone = $%d", len(args)))
	}

	if dto.Currency != nil {
		args = append(args, *dto.Currency)
		fields = append(fields, fmt.Sprintf("currency = $%d", len(args)))
	}

	if dto.Industry != nil {
		args = append(args, *dto.Industry)
		fields = append(fields, fmt.Sprintf("industry = $%d", len(args)))
	}

	if dto.Website != nil {
		args = append(args, *dto.Website)
		fields = append(fields, fmt.Sprintf("website = $%d", len(args)))
	}

	if dto.Description != nil {
		args = append(args, *dto.Description)
		fields = append(fields, fmt.Sprintf("description = $%d", len(args)))
	}

	query.WriteString(strings.Join(fields, ", "))

	args = append(args, workspaceID)

	fmt.Fprintf(
		&query,
		" WHERE id = $%d RETURNING id, name, slug, type, country, timezone, currency, industry, website, description, status, created_by, created_at, updated_at",
		len(args),
	)

	var workspace Workspace

	err := r.db.QueryRowContext(
		ctx,
		query.String(),
		args...,
	).Scan(
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
		return nil, err
	}

	return &workspace, nil
}

func (r *repository) Delete(
	ctx context.Context,
	workspaceID string,
) error {
	const query = `
        UPDATE workspaces
        SET
            status = 'deleted',
            updated_at = NOW()
        WHERE id = $1
    `

	result, err := r.db.ExecContext(ctx, query, workspaceID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
