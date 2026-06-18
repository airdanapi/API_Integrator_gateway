package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadUsesDefaultsAndRequiredDatabaseValues(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("BACKEND_PORT", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_HOST", "mysql")
	t.Setenv("DB_NAME", "api_integrator")
	t.Setenv("DB_USER", "gateway")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("JWT_SECRET", "01234567890123456789012345678901")
	t.Setenv("JWT_TTL", "")
	t.Setenv("JWT_ISSUER", "")
	t.Setenv("SEED_USERS_ENABLED", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned an unexpected error: %v", err)
	}

	if cfg.AppEnv != "development" {
		t.Fatalf("AppEnv = %q, want development", cfg.AppEnv)
	}
	if cfg.BackendPort != "8080" {
		t.Fatalf("BackendPort = %q, want 8080", cfg.BackendPort)
	}
	if cfg.DBPort != "3306" {
		t.Fatalf("DBPort = %q, want 3306", cfg.DBPort)
	}
	if cfg.DBHost != "mysql" || cfg.DBName != "api_integrator" || cfg.DBUser != "gateway" || cfg.DBPassword != "secret" {
		t.Fatalf("database configuration was not loaded correctly: %#v", cfg)
	}
	if cfg.JWTTTL != time.Hour {
		t.Fatalf("JWTTTL = %s, want 1h", cfg.JWTTTL)
	}
	if cfg.JWTIssuer != "api-integrator-gateway" {
		t.Fatalf("JWTIssuer = %q, want api-integrator-gateway", cfg.JWTIssuer)
	}
	if cfg.SeedUsersEnabled {
		t.Fatal("SeedUsersEnabled = true, want false")
	}
}

func TestLoadReportsEveryMissingRequiredVariable(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("JWT_SECRET", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}

	for _, name := range []string{"DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD", "JWT_SECRET"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("Load() error %q does not mention %s", err, name)
		}
	}
}

func TestLoadRejectsShortJWTSecretAndInvalidTTL(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("JWT_SECRET", "too-short")
	t.Setenv("JWT_TTL", "not-a-duration")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}
	if !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Errorf("Load() error %q does not mention JWT_SECRET", err)
	}
	if !strings.Contains(err.Error(), "JWT_TTL") {
		t.Errorf("Load() error %q does not mention JWT_TTL", err)
	}
}

func TestLoadRequiresSeedCredentialsWhenEnabled(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("SEED_USERS_ENABLED", "true")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want seed validation error")
	}

	for _, name := range []string{
		"SEED_ADMIN_USERNAME",
		"SEED_ADMIN_PASSWORD",
		"SEED_APP_USER_USERNAME",
		"SEED_APP_USER_PASSWORD",
		"SEED_APP_USER_APP_NAME",
		"SEED_MONITORING_USERNAME",
		"SEED_MONITORING_PASSWORD",
	} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("Load() error %q does not mention %s", err, name)
		}
	}
}

func TestLoadRejectsSeedUsersInProduction(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("APP_ENV", "production")
	t.Setenv("SEED_USERS_ENABLED", "true")
	setSeedEnvironment(t)

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want production seed validation error")
	}
	if !strings.Contains(err.Error(), "production") {
		t.Errorf("Load() error %q does not mention production", err)
	}
}

func setRequiredEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("APP_ENV", "test")
	t.Setenv("DB_HOST", "mysql")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_NAME", "api_integrator")
	t.Setenv("DB_USER", "gateway")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("JWT_SECRET", "01234567890123456789012345678901")
	t.Setenv("JWT_TTL", "1h")
	t.Setenv("JWT_ISSUER", "api-integrator-gateway")
	t.Setenv("SEED_USERS_ENABLED", "false")
}

func setSeedEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("SEED_ADMIN_USERNAME", "admin")
	t.Setenv("SEED_ADMIN_PASSWORD", "admin-password")
	t.Setenv("SEED_APP_USER_USERNAME", "marketplace")
	t.Setenv("SEED_APP_USER_PASSWORD", "marketplace-password")
	t.Setenv("SEED_APP_USER_APP_NAME", "Marketplace")
	t.Setenv("SEED_MONITORING_USERNAME", "insight")
	t.Setenv("SEED_MONITORING_PASSWORD", "insight-password")
}
