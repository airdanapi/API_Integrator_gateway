package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultAppEnv      = "development"
	defaultBackendPort = "8080"
	defaultDBPort      = "3306"
	defaultJWTIssuer   = "api-integrator-gateway"
	defaultJWTTTL      = time.Hour
	minimumJWTSecret   = 32
)

type SeedUserConfig struct {
	Username string
	Password string
	AppName  string
}

type Config struct {
	AppEnv             string
	BackendPort        string
	DBHost             string
	DBPort             string
	DBName             string
	DBUser             string
	DBPassword         string
	JWTSecret          string
	JWTTTL             time.Duration
	JWTIssuer          string
	SeedUsersEnabled   bool
	AdminSeed          SeedUserConfig
	AppUserSeed        SeedUserConfig
	MonitoringUserSeed SeedUserConfig
	GatewaySmartBankURL   string
	GatewayMarketplaceURL string
	GatewayLogisticsURL   string
	GatewaySupplierHubURL string
}

func Load() (Config, error) {
	jwtTTL, jwtTTLErr := durationOrDefault("JWT_TTL", defaultJWTTTL)
	seedUsersEnabled, seedEnabledErr := boolOrDefault("SEED_USERS_ENABLED", false)

	cfg := Config{
		AppEnv:           valueOrDefault("APP_ENV", defaultAppEnv),
		BackendPort:      valueOrDefault("BACKEND_PORT", defaultBackendPort),
		DBHost:           strings.TrimSpace(os.Getenv("DB_HOST")),
		DBPort:           valueOrDefault("DB_PORT", defaultDBPort),
		DBName:           strings.TrimSpace(os.Getenv("DB_NAME")),
		DBUser:           strings.TrimSpace(os.Getenv("DB_USER")),
		DBPassword:       strings.TrimSpace(os.Getenv("DB_PASSWORD")),
		JWTSecret:        strings.TrimSpace(os.Getenv("JWT_SECRET")),
		JWTTTL:           jwtTTL,
		JWTIssuer:        valueOrDefault("JWT_ISSUER", defaultJWTIssuer),
		SeedUsersEnabled: seedUsersEnabled,
		AdminSeed: SeedUserConfig{
			Username: strings.TrimSpace(os.Getenv("SEED_ADMIN_USERNAME")),
			Password: os.Getenv("SEED_ADMIN_PASSWORD"),
			AppName:  "API Gateway",
		},
		AppUserSeed: SeedUserConfig{
			Username: strings.TrimSpace(os.Getenv("SEED_APP_USER_USERNAME")),
			Password: os.Getenv("SEED_APP_USER_PASSWORD"),
			AppName:  strings.TrimSpace(os.Getenv("SEED_APP_USER_APP_NAME")),
		},
		MonitoringUserSeed: SeedUserConfig{
			Username: strings.TrimSpace(os.Getenv("SEED_MONITORING_USERNAME")),
			Password: os.Getenv("SEED_MONITORING_PASSWORD"),
			AppName:  "UMKM Insight",
		},
		GatewaySmartBankURL:   strings.TrimSpace(os.Getenv("GATEWAY_SMARTBANK_URL")),
		GatewayMarketplaceURL: strings.TrimSpace(os.Getenv("GATEWAY_MARKETPLACE_URL")),
		GatewayLogisticsURL:   strings.TrimSpace(os.Getenv("GATEWAY_LOGISTICS_URL")),
		GatewaySupplierHubURL: strings.TrimSpace(os.Getenv("GATEWAY_SUPPLIERHUB_URL")),
	}

	validationErrors := make([]string, 0, 8)
	for name, value := range map[string]string{
		"DB_HOST":     cfg.DBHost,
		"DB_NAME":     cfg.DBName,
		"DB_USER":     cfg.DBUser,
		"DB_PASSWORD": cfg.DBPassword,
		"JWT_SECRET":  cfg.JWTSecret,
	} {
		if value == "" {
			validationErrors = append(validationErrors, name+" is required")
		}
	}
	if cfg.JWTSecret != "" && len(cfg.JWTSecret) < minimumJWTSecret {
		validationErrors = append(
			validationErrors,
			fmt.Sprintf("JWT_SECRET must contain at least %d characters", minimumJWTSecret),
		)
	}
	if jwtTTLErr != nil {
		validationErrors = append(validationErrors, jwtTTLErr.Error())
	}
	if seedEnabledErr != nil {
		validationErrors = append(validationErrors, seedEnabledErr.Error())
	}
	if cfg.SeedUsersEnabled {
		if strings.EqualFold(cfg.AppEnv, "production") {
			validationErrors = append(
				validationErrors,
				"SEED_USERS_ENABLED cannot be true in production",
			)
		}
		for name, value := range map[string]string{
			"SEED_ADMIN_USERNAME":      cfg.AdminSeed.Username,
			"SEED_ADMIN_PASSWORD":      cfg.AdminSeed.Password,
			"SEED_APP_USER_USERNAME":   cfg.AppUserSeed.Username,
			"SEED_APP_USER_PASSWORD":   cfg.AppUserSeed.Password,
			"SEED_APP_USER_APP_NAME":   cfg.AppUserSeed.AppName,
			"SEED_MONITORING_USERNAME": cfg.MonitoringUserSeed.Username,
			"SEED_MONITORING_PASSWORD": cfg.MonitoringUserSeed.Password,
		} {
			if strings.TrimSpace(value) == "" {
				validationErrors = append(validationErrors, name+" is required when seeding")
			}
		}
	}
	if len(validationErrors) > 0 {
		return Config{}, fmt.Errorf(
			"invalid configuration: %s",
			strings.Join(validationErrors, "; "),
		)
	}

	return cfg, nil
}

func valueOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func durationOrDefault(name string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback, nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration", name)
	}
	return duration, nil
}

func boolOrDefault(name string, fallback bool) (bool, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean", name)
	}
	return parsed, nil
}
