package workspace

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/thakur-vijay/adorbit-go/apperrors"
)

type mockRepository struct {
	mu sync.Mutex

	createWorkspace *Workspace
	createErr       error
	fetchWorkspaces []Workspace
	fetchErr        error

	createCalls int
	fetchCalls  int
	createDTO   CreateWorkspaceDTO
	createUser  string
	createSlug  string
	fetchUser   string
	createCtx   context.Context
	fetchCtx    context.Context
}

func (m *mockRepository) Create(ctx context.Context, dto CreateWorkspaceDTO, userID, slug string) (*Workspace, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createCalls++
	m.createCtx = ctx
	m.createDTO = dto
	m.createUser = userID
	m.createSlug = slug
	return m.createWorkspace, m.createErr
}

func (m *mockRepository) Fetch(ctx context.Context, userID string) ([]Workspace, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.fetchCalls++
	m.fetchCtx = ctx
	m.fetchUser = userID
	return m.fetchWorkspaces, m.fetchErr
}

func TestCreateService_GeneratesSlugAndForwardsInput(t *testing.T) {
	tests := []struct {
		name string
		slug string
	}{
		{name: "My Workspace", slug: "my-workspace"},
		{name: "AdOrbit Publisher", slug: "adorbit-publisher"},
		{name: "ADORBIT Publisher", slug: "adorbit-publisher"},
		{name: "  Alpha Workspace  ", slug: "alpha-workspace"},
		{name: "Alpha   Beta", slug: "alpha-beta"},
		{name: "Café", slug: "cafe"},
		{name: "Alpha / Beta", slug: "alpha-beta"},
		{name: "My Workspace 2026", slug: "my-workspace-2026"},
		{name: "!!!", slug: ""},
		{name: "", slug: ""},
		{name: "   ", slug: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dto := testWorkspaceDTO()
			dto.Name = test.name
			originalDTO := dto
			userID := uuid.NewString()
			ctx := context.WithValue(context.Background(), struct{}{}, "request-context")
			wantWorkspace := testWorkspace(userID, test.name)
			repository := &mockRepository{createWorkspace: wantWorkspace}
			service := NewService(repository)

			got, err := service.Create(ctx, userID, dto)
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}
			if !reflect.DeepEqual(got, wantWorkspace) {
				t.Errorf("Create() workspace = %#v, want %#v", got, wantWorkspace)
			}

			repository.mu.Lock()
			defer repository.mu.Unlock()
			if repository.createCalls != 1 {
				t.Errorf("repository Create calls = %d, want 1", repository.createCalls)
			}
			if repository.createUser != userID {
				t.Errorf("repository UserID = %q, want %q", repository.createUser, userID)
			}
			if repository.createSlug != test.slug {
				t.Errorf("repository slug = %q, want %q", repository.createSlug, test.slug)
			}
			if !reflect.DeepEqual(repository.createDTO, originalDTO) {
				t.Errorf("repository DTO = %#v, want %#v", repository.createDTO, originalDTO)
			}
			if repository.createCtx != ctx {
				t.Errorf("repository context = %v, want original context", repository.createCtx)
			}
			if !reflect.DeepEqual(dto, originalDTO) {
				t.Errorf("Create() modified DTO: got %#v, want %#v", dto, originalDTO)
			}
		})
	}
}

func TestCreateService_ErrorClassificationAndNilResults(t *testing.T) {
	conflict := apperrors.NewConflict("duplicate workspace", nil)
	uniqueViolation := &pgconn.PgError{Code: "23505", Message: "duplicate key"}
	tests := []struct {
		name          string
		result        *Workspace
		repositoryErr error
		wantNil       bool
		wantStatus    int
		wantIs        error
	}{
		{
			name:          "repository error is classified as internal",
			repositoryErr: errors.New("database failure"),
			wantNil:       true,
			wantStatus:    500,
		},
		{
			name:          "repository context cancellation retains timeout semantics",
			repositoryErr: context.Canceled,
			wantNil:       true,
			wantStatus:    408,
			wantIs:        context.Canceled,
		},
		{
			name:          "repository conflict classification is retained",
			repositoryErr: conflict,
			wantNil:       true,
			wantStatus:    409,
			wantIs:        conflict,
		},
		{
			name:          "PostgreSQL unique violation is classified as conflict",
			repositoryErr: uniqueViolation,
			wantNil:       true,
			wantStatus:    409,
			wantIs:        uniqueViolation,
		},
		{
			name:    "nil repository result with nil error remains nil",
			wantNil: true,
		},
		{
			name:          "workspace and repository error returns classified error",
			result:        testWorkspace(uuid.NewString(), "partial result"),
			repositoryErr: errors.New("write failed"),
			wantNil:       true,
			wantStatus:    500,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &mockRepository{
				createWorkspace: test.result,
				createErr:       test.repositoryErr,
			}
			service := NewService(repository)
			got, err := service.Create(context.Background(), uuid.NewString(), testWorkspaceDTO())

			if test.repositoryErr == nil {
				if err != nil {
					t.Fatalf("Create() error = %v, want nil", err)
				}
			} else {
				var appErr *apperrors.AppError
				if !errors.As(err, &appErr) {
					t.Fatalf("Create() error = %v, want *apperrors.AppError", err)
				}
				if appErr.StatusCode != test.wantStatus {
					t.Errorf("classified status = %d, want %d", appErr.StatusCode, test.wantStatus)
				}
				if test.wantIs != nil && !errors.Is(err, test.wantIs) {
					t.Errorf("Create() error %v does not wrap %v", err, test.wantIs)
				}
			}
			if test.wantNil && got != nil {
				t.Errorf("Create() workspace = %#v, want nil", got)
			}

			repository.mu.Lock()
			defer repository.mu.Unlock()
			if repository.createCalls != 1 {
				t.Errorf("repository Create calls = %d, want 1", repository.createCalls)
			}
		})
	}
}

func TestFetchService_ResultsAndErrors(t *testing.T) {
	userID := uuid.NewString()
	first := *testWorkspace(userID, "First")
	second := *testWorkspace(userID, "Second")
	second.Industry = nil
	second.Website = nil
	second.Description = nil
	repositoryError := errors.New("fetch failed")

	tests := []struct {
		name          string
		result        []Workspace
		repositoryErr error
		wantResult    []Workspace
		wantStatus    int
		wantIs        error
	}{
		{name: "multiple workspaces preserve values and order", result: []Workspace{first, second}, wantResult: []Workspace{first, second}},
		{name: "empty slice remains non-nil", result: []Workspace{}, wantResult: []Workspace{}},
		{name: "nil slice remains nil", result: nil, wantResult: nil},
		{name: "repository error is classified", repositoryErr: repositoryError, wantStatus: 500, wantIs: repositoryError},
		{name: "repository cancellation preserves timeout semantics", repositoryErr: context.Canceled, wantStatus: 408, wantIs: context.Canceled},
		{name: "repository deadline preserves timeout semantics", repositoryErr: context.DeadlineExceeded, wantStatus: 408, wantIs: context.DeadlineExceeded},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), struct{}{}, "fetch-context")
			repository := &mockRepository{
				fetchWorkspaces: test.result,
				fetchErr:        test.repositoryErr,
			}
			service := NewService(repository)
			got, err := service.Fetch(ctx, userID)

			if test.repositoryErr == nil {
				if err != nil {
					t.Fatalf("Fetch() error = %v, want nil", err)
				}
				if !reflect.DeepEqual(got, test.wantResult) {
					t.Errorf("Fetch() = %#v, want %#v", got, test.wantResult)
				}
			} else {
				var appErr *apperrors.AppError
				if !errors.As(err, &appErr) {
					t.Fatalf("Fetch() error = %v, want *apperrors.AppError", err)
				}
				if appErr.StatusCode != test.wantStatus {
					t.Errorf("classified status = %d, want %d", appErr.StatusCode, test.wantStatus)
				}
				if !errors.Is(err, test.wantIs) {
					t.Errorf("Fetch() error %v does not wrap %v", err, test.wantIs)
				}
				if got != nil {
					t.Errorf("Fetch() result = %#v, want nil on error", got)
				}
			}

			repository.mu.Lock()
			defer repository.mu.Unlock()
			if repository.fetchCalls != 1 {
				t.Errorf("repository Fetch calls = %d, want 1", repository.fetchCalls)
			}
			if repository.fetchUser != userID {
				t.Errorf("repository UserID = %q, want %q", repository.fetchUser, userID)
			}
			if repository.fetchCtx != ctx {
				t.Errorf("repository context = %v, want original context", repository.fetchCtx)
			}
		})
	}
}
