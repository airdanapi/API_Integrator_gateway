package config

import (
	"strings"
	"testing"
)

func TestLoadUsesDefaultsAndRequiredDatabaseValues(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("BACKEND_PORT", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_HOST", "mysql")
	t.Setenv("DB_NAME", "api_integrator")
	t.Setenv("DB_USER", "gateway")
	t.Setenv("DB_PASSWORD", "secret")

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
}

func TestLoadReportsEveryMissingRequiredDatabaseVariable(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}

	for _, name := range []string{"DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("Load() error %q does not mention %s", err, name)
		}
	}
}
