package workspace

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/thakur-vijay/adorbit-go/idempotency"
	"github.com/thakur-vijay/adorbit-go/middleware"
)

const testJWTIssuer = "workspace-tests"

type jwtTestFixture struct {
	privateKey *rsa.PrivateKey
	publicPEM  string
}

func newJWTTestFixture(t *testing.T) *jwtTestFixture {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate test JWT key: %v", err)
	}
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal test JWT public key: %v", err)
	}

	return &jwtTestFixture{
		privateKey: privateKey,
		publicPEM: string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyDER,
		})),
	}
}

func (f *jwtTestFixture) middleware(t *testing.T) *middleware.JWTMiddleware {
	t.Helper()

	auth, err := middleware.NewJWTMiddleware(f.publicPEM, testJWTIssuer)
	if err != nil {
		t.Fatalf("create test JWT middleware: %v", err)
	}
	return auth
}

func (f *jwtTestFixture) token(t *testing.T, userID string, expiresAt time.Time) string {
	t.Helper()

	now := time.Now()
	claims := middleware.JWTClaims{
		UserID: userID,
		Role:   "advertiser",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    testJWTIssuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(f.privateKey)
	if err != nil {
		t.Fatalf("sign test JWT: %v", err)
	}
	return token
}

func testWorkspaceDTO() CreateWorkspaceDTO {
	industry := "Retail"
	website := "https://example.test"
	description := "Primary workspace"
	return CreateWorkspaceDTO{
		Name:        "AdOrbit Publisher",
		Type:        Publisher,
		Country:     "US",
		Timezone:    "America/New_York",
		Currency:    "USD",
		Industry:    &industry,
		Website:     &website,
		Description: &description,
	}
}

func stringPointer(value string) *string { return &value }

func testWorkspace(userID, name string) *Workspace {
	now := time.Now().UTC().Truncate(time.Microsecond)
	return &Workspace{
		ID:          uuid.NewString(),
		Name:        name,
		Slug:        strings.ToLower(strings.ReplaceAll(name, " ", "-")),
		Type:        Publisher,
		Country:     "US",
		Timezone:    "America/New_York",
		Currency:    "USD",
		Status:      Active,
		CreatedBy:   userID,
		CreatedAt:   now,
		UpdatedAt:   now,
		Industry:    stringPointer("Retail"),
		Website:     stringPointer("https://example.test"),
		Description: stringPointer("Test workspace"),
	}
}

// WORKSPACE_TEST_DATABASE_URL must point to a disposable PostgreSQL database.
// Each test gets its own schema, and cleanup drops only that generated schema.
func workspaceTestDB(t *testing.T, withIdempotency bool) *sql.DB {
	t.Helper()

	dsn := os.Getenv("WORKSPACE_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("WORKSPACE_TEST_DATABASE_URL is not set")
	}

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("parse test database URL: %v", err)
	}

	adminDB := stdlib.OpenDB(*config)
	t.Cleanup(func() { _ = adminDB.Close() })

	schema := "workspace_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := adminDB.ExecContext(context.Background(), "CREATE SCHEMA "+schema); err != nil {
		t.Fatalf("create isolated test schema: %v", err)
	}

	testConfig := config.Copy()
	if testConfig.RuntimeParams == nil {
		testConfig.RuntimeParams = make(map[string]string)
	}
	testConfig.RuntimeParams["search_path"] = schema
	db := stdlib.OpenDB(*testConfig)
	db.SetMaxOpenConns(12)
	t.Cleanup(func() {
		_ = db.Close()
		if _, err := adminDB.ExecContext(context.Background(), "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Errorf("drop isolated test schema: %v", err)
		}
	})

	migration, err := os.ReadFile("../database/migrations/000001_create_workspace.up.sql")
	if err != nil {
		t.Fatalf("read workspace migration: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), string(migration)); err != nil {
		t.Fatalf("apply workspace migration: %v", err)
	}

	if withIdempotency {
		migrations := []string{
			"migrations/001_create_idempotency_keys.up.sql",
			"migrations/002_add_processing_owner.up.sql",
			"migrations/003_keyrings_and_indeterminate.up.sql",
		}
		for _, filename := range migrations {
			migration, err := idempotency.Migrations.ReadFile(filename)
			if err != nil {
				t.Fatalf("read idempotency migration %s: %v", filename, err)
			}
			if _, err := db.ExecContext(context.Background(), string(migration)); err != nil {
				t.Fatalf("apply idempotency migration %s: %v", filename, err)
			}
		}
	}

	return db
}

func newIdempotencyMiddleware(t *testing.T, db *sql.DB) *idempotency.Middleware {
	t.Helper()

	const scopeKeyID = "scope-v1"
	const responseKeyID = "response-v1"
	scopeKey := []byte(strings.Repeat("s", 32))
	responseKey := []byte(strings.Repeat("e", 32))
	service, err := idempotency.NewService(
		idempotency.NewPostgresRepository(db),
		idempotency.Config{
			TTL:           time.Hour,
			ProcessingTTL: time.Second,
			ScopeHashKeys: []idempotency.KeyMaterial{
				{ID: scopeKeyID, Key: scopeKey},
			},
			ActiveScopeHashKeyID: scopeKeyID,
			ResponseEncryptionKeys: []idempotency.KeyMaterial{
				{ID: responseKeyID, Key: responseKey},
			},
			ActiveResponseKeyID: responseKeyID,
		},
	)
	if err != nil {
		t.Fatalf("create test idempotency service: %v", err)
	}
	return idempotency.NewHTTPMiddleware(service)
}
