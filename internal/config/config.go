package config

import (
	"encoding/hex"
	"errors"
	"os"
	"strings"

	"github.com/thakur-vijay/adorbit-go/idempotency"
)

type Config struct {
	Port                     string
	DatabaseURL              string
	JWTPublicKey             string
	IdempotencyScopeKey      []byte
	IdempotencyEncryptionKey []byte
	SWAGGER_HOST             string
}

func Load() (*Config, error) {
	port := normalizePort(getEnv("APP_PORT", "8080"))

	scopeKey, err := hex.DecodeString(os.Getenv("IDEMPOTENCY_SCOPE_KEY"))
	if err != nil || len(scopeKey) != 32 {
		return nil, errors.New("IDEMPOTENCY_SCOPE_KEY must be 64 hexadecimal characters")
	}

	encryptionKey, err := hex.DecodeString(os.Getenv("IDEMPOTENCY_ENCRYPTION_KEY"))
	if err != nil || len(encryptionKey) != 32 {
		return nil, errors.New("IDEMPOTENCY_ENCRYPTION_KEY must be 64 hexadecimal characters")
	}

	cfg := &Config{
		Port:                     port,
		DatabaseURL:              getEnv("WORKSPACE_DATABASE_URL", ""),
		JWTPublicKey:             strings.ReplaceAll(os.Getenv("JWT_PUBLIC_KEY"), `\n`, "\n"),
		IdempotencyScopeKey:      scopeKey,
		IdempotencyEncryptionKey: encryptionKey,
		SWAGGER_HOST:             getEnv("SWAGGER_HOST", ""),
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("WORKSPACE_DATABASE_URL is not set")
	}

	if cfg.JWTPublicKey == "" {
		return nil, errors.New("JWT_PUBLIC_KEY is not set")
	}

	return cfg, nil
}

func (c *Config) IdempotencyConfig() idempotency.Config {
	return idempotency.Config{
		TTL:           24 * 60 * 60 * 1e9,
		ProcessingTTL: 2 * 60 * 1e9,

		ScopeHashKeys: []idempotency.KeyMaterial{
			{
				ID:  "scope-v1",
				Key: c.IdempotencyScopeKey,
			},
		},
		ActiveScopeHashKeyID: "scope-v1",

		ResponseEncryptionKeys: []idempotency.KeyMaterial{
			{
				ID:  "response-v1",
				Key: c.IdempotencyEncryptionKey,
			},
		},
		ActiveResponseKeyID: "response-v1",
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func normalizePort(port string) string {
	if port == "" {
		return ":8080"
	}

	if strings.Contains(port, ":") {
		return port
	}

	return ":" + port
}
