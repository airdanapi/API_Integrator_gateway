package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	defaultAppEnv      = "development"
	defaultBackendPort = "8080"
	defaultDBPort      = "3306"
)

type Config struct {
	AppEnv      string
	BackendPort string
	DBHost      string
	DBPort      string
	DBName      string
	DBUser      string
	DBPassword  string
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:      valueOrDefault("APP_ENV", defaultAppEnv),
		BackendPort: valueOrDefault("BACKEND_PORT", defaultBackendPort),
		DBHost:      strings.TrimSpace(os.Getenv("DB_HOST")),
		DBPort:      valueOrDefault("DB_PORT", defaultDBPort),
		DBName:      strings.TrimSpace(os.Getenv("DB_NAME")),
		DBUser:      strings.TrimSpace(os.Getenv("DB_USER")),
		DBPassword:  strings.TrimSpace(os.Getenv("DB_PASSWORD")),
	}

	missing := make([]string, 0, 4)
	for name, value := range map[string]string{
		"DB_HOST":     cfg.DBHost,
		"DB_NAME":     cfg.DBName,
		"DB_USER":     cfg.DBUser,
		"DB_PASSWORD": cfg.DBPassword,
	} {
		if value == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func valueOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}
