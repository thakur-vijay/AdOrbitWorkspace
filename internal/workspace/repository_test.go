package workspace

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func repositoryDTO(name string, workspaceType WorkspaceType) CreateWorkspaceDTO {
	return CreateWorkspaceDTO{
		Name:     name,
		Type:     workspaceType,
		Country:  "US",
		Timezone: "America/New_York",
		Currency: "USD",
	}
}

func createRepositoryWorkspace(t *testing.T, repository Repository, userID, name string, workspaceType WorkspaceType, slug string, dto CreateWorkspaceDTO) *Workspace {
	t.Helper()
	if slug == "" {
		slug = fmt.Sprintf("%s-%s", name, uuid.NewString())
	}
	created, err := repository.Create(context.Background(), dto, userID, slug)
	if err != nil {
		t.Fatalf("create workspace %q: %v", name, err)
	}
	return created
}

func TestCreateRepository_PostgreSQL(t *testing.T) {
	db := workspaceTestDB(t, false)
	repository := NewRepository(db)
	userID := uuid.NewString()

	populated := repositoryDTO("Advertiser workspace", Advertiser)
	populated.Industry = stringPointer("Retail")
	populated.Website = stringPointer("https://advertiser.example.test")
	populated.Description = stringPointer("A fully populated workspace")
	nullable := repositoryDTO("Publisher workspace", Publisher)

	typeCase := []struct {
		name         string
		dto          CreateWorkspaceDTO
		slug         string
		wantIndustry *string
		wantWebsite  *string
		wantDesc     *string
	}{
		{
			name:         "advertiser with nullable fields populated",
			dto:          populated,
			slug:         "advertiser-workspace",
			wantIndustry: populated.Industry,
			wantWebsite:  populated.Website,
			wantDesc:     populated.Description,
		},
		{
			name: "publisher with nullable fields nil",
			dto:  nullable,
			slug: "publisher-workspace",
		},
	}

	ids := make(map[string]bool)
	for _, test := range typeCase {
		t.Run(test.name, func(t *testing.T) {
			started := time.Now().Add(-time.Second)
			got, err := repository.Create(context.Background(), test.dto, userID, test.slug)
			finished := time.Now().Add(time.Second)
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}
			if _, err := uuid.Parse(got.ID); err != nil {
				t.Errorf("returned ID %q is not a UUID: %v", got.ID, err)
			}
			if ids[got.ID] {
				t.Errorf("generated duplicate UUID %q", got.ID)
			}
			ids[got.ID] = true
			if got.Name != test.dto.Name || got.Slug != test.slug || got.Type != test.dto.Type {
				t.Errorf("returned identity fields = (%q, %q, %q), want (%q, %q, %q)", got.Name, got.Slug, got.Type, test.dto.Name, test.slug, test.dto.Type)
			}
			if got.Country != test.dto.Country || got.Timezone != test.dto.Timezone || got.Currency != test.dto.Currency || got.CreatedBy != userID {
				t.Errorf("returned workspace fields do not match inserted DTO/user: %#v", got)
			}
			if got.Status != Active {
				t.Errorf("status = %q, want %q", got.Status, Active)
			}
			if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() || got.CreatedAt.Before(started) || got.CreatedAt.After(finished) || got.UpdatedAt.Before(started) || got.UpdatedAt.After(finished) {
				t.Errorf("timestamps not populated in request window: created_at=%v updated_at=%v", got.CreatedAt, got.UpdatedAt)
			}
			if got.Industry == nil && test.wantIndustry != nil || got.Industry != nil && test.wantIndustry == nil || got.Industry != nil && *got.Industry != *test.wantIndustry {
				t.Errorf("industry = %v, want %v", got.Industry, test.wantIndustry)
			}
			if got.Website == nil && test.wantWebsite != nil || got.Website != nil && test.wantWebsite == nil || got.Website != nil && *got.Website != *test.wantWebsite {
				t.Errorf("website = %v, want %v", got.Website, test.wantWebsite)
			}
			if got.Description == nil && test.wantDesc != nil || got.Description != nil && test.wantDesc == nil || got.Description != nil && *got.Description != *test.wantDesc {
				t.Errorf("description = %v, want %v", got.Description, test.wantDesc)
			}
		})
	}
	if len(ids) != 2 {
		t.Errorf("unique generated IDs = %d, want 2", len(ids))
	}

	secondUserWorkspace := createRepositoryWorkspace(t, repository, uuid.NewString(), "Second record", Advertiser, "", repositoryDTO("Second record", Advertiser))
	if ids[secondUserWorkspace.ID] {
		t.Errorf("generated UUID %q was reused", secondUserWorkspace.ID)
	}
}

func TestCreateRepository_ConstraintFailures(t *testing.T) {
	db := workspaceTestDB(t, false)
	repository := NewRepository(db)
	userID := uuid.NewString()

	tests := []struct {
		name   string
		dto    CreateWorkspaceDTO
		userID string
		pgCode string
	}{
		{
			name: "workspace type check constraint",
			dto: func() CreateWorkspaceDTO {
				dto := repositoryDTO("invalid type", Publisher)
				dto.Type = "agency"
				return dto
			}(),
			userID: userID,
			pgCode: "23514",
		},
		{
			name: "country exceeds char width",
			dto: func() CreateWorkspaceDTO {
				dto := repositoryDTO("invalid country", Publisher)
				dto.Country = "USA"
				return dto
			}(),
			userID: userID,
			pgCode: "22001",
		},
		{
			name: "currency exceeds char width",
			dto: func() CreateWorkspaceDTO {
				dto := repositoryDTO("invalid currency", Publisher)
				dto.Currency = "USDD"
				return dto
			}(),
			userID: userID,
			pgCode: "22001",
		},
		{
			name:   "created by is not a UUID",
			dto:    repositoryDTO("invalid user ID", Publisher),
			userID: "not-a-uuid",
			pgCode: "22P02",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, err := repository.Create(context.Background(), test.dto, test.userID, uuid.NewString()); err == nil {
				t.Fatalf("Create() = %#v, want database constraint error", got)
			} else {
				var pgErr *pgconn.PgError
				if !errors.As(err, &pgErr) {
					t.Fatalf("Create() error = %T %v, want *pgconn.PgError", err, err)
				}
				if pgErr.Code != test.pgCode {
					t.Errorf("PostgreSQL error code = %q, want %q", pgErr.Code, test.pgCode)
				}
			}
		})
	}

	// PostgreSQL CHAR(n) pads shorter values; the migration does not require
	// exact country/currency input lengths, so short values are accepted.
	shortFields := repositoryDTO("short fields", Publisher)
	shortFields.Country = "U"
	shortFields.Currency = "U"
	if _, err := repository.Create(context.Background(), shortFields, userID, "short-field-widths"); err != nil {
		t.Fatalf("short CHAR fields unexpectedly rejected: %v", err)
	}
}

func TestCreateRepository_DuplicateSlugDoesNotInsertDuplicate(t *testing.T) {
	db := workspaceTestDB(t, false)
	repository := NewRepository(db)
	userID := uuid.NewString()
	dto := repositoryDTO("Duplicate slug", Publisher)
	slug := "shared-unique-slug"

	if _, err := repository.Create(context.Background(), dto, userID, slug); err != nil {
		t.Fatalf("first Create() error = %v", err)
	}
	if _, err := repository.Create(context.Background(), dto, userID, slug); err == nil {
		t.Fatal("second Create() with duplicate slug succeeded")
	} else {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
			t.Errorf("duplicate slug error = %T %v, want unique constraint violation 23505", err, err)
		}
	}

	var count int
	if err := db.QueryRowContext(context.Background(), "SELECT count(*) FROM workspaces WHERE slug = $1", slug).Scan(&count); err != nil {
		t.Fatalf("count duplicate-slug rows: %v", err)
	}
	if count != 1 {
		t.Errorf("stored rows for duplicate slug = %d, want 1", count)
	}
}

func TestCreateRepository_DatabaseAndContextErrors(t *testing.T) {
	t.Run("database unavailable", func(t *testing.T) {
		db := workspaceTestDB(t, false)
		repository := NewRepository(db)
		if err := db.Close(); err != nil {
			t.Fatalf("close test database: %v", err)
		}
		if _, err := repository.Create(context.Background(), repositoryDTO("closed DB", Publisher), uuid.NewString(), "closed-db"); err == nil {
			t.Fatal("Create() with closed database returned nil error")
		}
	})

	t.Run("context already cancelled", func(t *testing.T) {
		db := workspaceTestDB(t, false)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := NewRepository(db).Create(ctx, repositoryDTO("cancelled", Publisher), uuid.NewString(), "cancelled"); !errors.Is(err, context.Canceled) {
			t.Errorf("Create() error = %v, want context.Canceled", err)
		}
	})

	t.Run("context deadline exceeded", func(t *testing.T) {
		db := workspaceTestDB(t, false)
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
		defer cancel()
		if _, err := NewRepository(db).Create(ctx, repositoryDTO("expired", Publisher), uuid.NewString(), "expired"); !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("Create() error = %v, want context.DeadlineExceeded", err)
		}
	})
}

func TestFetchRepository_PostgreSQLIsolationOrderingAndNullableFields(t *testing.T) {
	db := workspaceTestDB(t, false)
	repository := NewRepository(db)
	userA, userB := uuid.NewString(), uuid.NewString()

	dto := repositoryDTO("Same name", Advertiser)
	first := createRepositoryWorkspace(t, repository, userA, dto.Name, Advertiser, "user-a-old", dto)
	secondDTO := repositoryDTO("Newest workspace", Publisher)
	second := createRepositoryWorkspace(t, repository, userA, secondDTO.Name, Publisher, "user-a-new", secondDTO)
	otherUser := createRepositoryWorkspace(t, repository, userB, dto.Name, Advertiser, "user-b-same-name", dto)
	deleted := createRepositoryWorkspace(t, repository, userA, "Deleted workspace", Publisher, "user-a-deleted", repositoryDTO("Deleted workspace", Publisher))

	oldTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	newTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, err := db.ExecContext(context.Background(), "UPDATE workspaces SET created_at=$1, updated_at=$1 WHERE id=$2", oldTime, first.ID); err != nil {
		t.Fatalf("set old deterministic timestamp: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), "UPDATE workspaces SET created_at=$1, updated_at=$1 WHERE id=$2", newTime, second.ID); err != nil {
		t.Fatalf("set new deterministic timestamp: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), "UPDATE workspaces SET created_at=$1, updated_at=$1 WHERE id=$2", newTime.Add(time.Hour), otherUser.ID); err != nil {
		t.Fatalf("set other user's timestamp: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), "UPDATE workspaces SET status='deleted' WHERE id=$1", deleted.ID); err != nil {
		t.Fatalf("mark workspace deleted: %v", err)
	}

	got, err := repository.Fetch(context.Background(), userA)
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("Fetch() returned %d workspaces, want 2: %#v", len(got), got)
	}
	if got[0].ID != second.ID || got[1].ID != first.ID {
		t.Errorf("Fetch() order = [%s, %s], want newest first [%s, %s]", got[0].ID, got[1].ID, second.ID, first.ID)
	}
	for _, workspace := range got {
		if workspace.CreatedBy != userA {
			t.Errorf("workspace %q belongs to %q, want only user %q", workspace.ID, workspace.CreatedBy, userA)
		}
		if workspace.Status != Active {
			t.Errorf("workspace %q status = %q, want active", workspace.ID, workspace.Status)
		}
		if workspace.Industry != nil || workspace.Website != nil || workspace.Description != nil {
			t.Errorf("nullable fields were not scanned as nil: %#v", workspace)
		}
	}

	deletedOnlyUser := uuid.NewString()
	deletedOnly := createRepositoryWorkspace(t, repository, deletedOnlyUser, "Only deleted", Advertiser, "only-deleted", repositoryDTO("Only deleted", Advertiser))
	if _, err := db.ExecContext(context.Background(), "UPDATE workspaces SET status='deleted' WHERE id=$1", deletedOnly.ID); err != nil {
		t.Fatalf("mark only workspace deleted: %v", err)
	}
	empty, err := repository.Fetch(context.Background(), deletedOnlyUser)
	if err != nil {
		t.Fatalf("Fetch() for deleted-only user error = %v", err)
	}
	if empty == nil || len(empty) != 0 {
		t.Errorf("Fetch() for deleted-only user = %#v, want non-nil empty slice", empty)
	}
}

func TestFetchRepository_QueryScanAndIterationErrors(t *testing.T) {
	t.Run("query error is propagated", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		queryErr := errors.New("query failed")
		mock.ExpectQuery("FROM workspaces").WithArgs("user-id").WillReturnError(queryErr)

		_, err = NewRepository(db).Fetch(context.Background(), "user-id")
		if !errors.Is(err, queryErr) {
			t.Errorf("Fetch() error = %v, want wrapped query error", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet SQL expectations: %v", err)
		}
	})

	t.Run("row scan error is propagated", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		rows := sqlmock.NewRows(workspaceColumns()).AddRow(
			nil, "name", "slug", "advertiser", "US", "UTC", "USD", nil, nil, nil,
			"active", uuid.NewString(), time.Now(), time.Now(),
		)
		mock.ExpectQuery("FROM workspaces").WithArgs("user-id").WillReturnRows(rows)

		if _, err := NewRepository(db).Fetch(context.Background(), "user-id"); err == nil {
			t.Fatal("Fetch() with a NULL workspace ID returned nil error")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet SQL expectations: %v", err)
		}
	})

	t.Run("rows iteration error is propagated", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		iterationErr := errors.New("iteration failed")
		rows := sqlmock.NewRows(workspaceColumns()).
			AddRow(uuid.NewString(), "name", "slug", "advertiser", "US", "UTC", "USD", nil, nil, nil, "active", uuid.NewString(), time.Now(), time.Now()).
			AddRow(uuid.NewString(), "second", "second-slug", "publisher", "US", "UTC", "USD", nil, nil, nil, "active", uuid.NewString(), time.Now(), time.Now()).
			RowError(1, iterationErr)
		mock.ExpectQuery("FROM workspaces").WithArgs("user-id").WillReturnRows(rows)

		_, err = NewRepository(db).Fetch(context.Background(), "user-id")
		if !errors.Is(err, iterationErr) {
			t.Errorf("Fetch() error = %v, want wrapped iteration error", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet SQL expectations: %v", err)
		}
	})
}

func TestFetchRepository_ContextErrors(t *testing.T) {
	db := workspaceTestDB(t, false)
	repository := NewRepository(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := repository.Fetch(ctx, uuid.NewString()); !errors.Is(err, context.Canceled) {
		t.Errorf("Fetch() with cancelled context error = %v, want context.Canceled", err)
	}

	expired, expireCancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer expireCancel()
	if _, err := repository.Fetch(expired, uuid.NewString()); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Fetch() with expired deadline error = %v, want context.DeadlineExceeded", err)
	}
}

func workspaceColumns() []string {
	return []string{
		"id", "name", "slug", "type", "country", "timezone", "currency",
		"industry", "website", "description", "status", "created_by",
		"created_at", "updated_at",
	}
}
