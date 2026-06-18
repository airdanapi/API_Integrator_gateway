package auth

import (
	"testing"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

func TestJWTServiceGeneratesAndValidatesRequiredClaims(t *testing.T) {
	now := time.Date(2026, 6, 18, 4, 0, 0, 0, time.UTC)
	service := NewJWTService(
		"01234567890123456789012345678901",
		"api-integrator-gateway",
		time.Hour,
		func() time.Time { return now },
	)
	user := model.User{
		ID:       42,
		Username: "admin",
		Role:     model.RoleAdminGateway,
		AppName:  "API Gateway",
	}

	token, expiresIn, err := service.Generate(user)
	if err != nil {
		t.Fatalf("Generate() returned an unexpected error: %v", err)
	}
	if expiresIn != 3600 {
		t.Fatalf("expiresIn = %d, want 3600", expiresIn)
	}

	claims, err := service.Validate(token)
	if err != nil {
		t.Fatalf("Validate() returned an unexpected error: %v", err)
	}
	if claims.Subject != "42" ||
		claims.Username != "admin" ||
		claims.Role != model.RoleAdminGateway ||
		claims.AppName != "API Gateway" ||
		claims.Issuer != "api-integrator-gateway" {
		t.Fatalf("unexpected claims: %#v", claims)
	}
	if claims.IssuedAt == nil || claims.NotBefore == nil || claims.ExpiresAt == nil {
		t.Fatalf("registered time claims must be present: %#v", claims)
	}
	if !claims.IssuedAt.Time.Equal(now) ||
		!claims.NotBefore.Time.Equal(now) ||
		!claims.ExpiresAt.Time.Equal(now.Add(time.Hour)) {
		t.Fatalf("unexpected token times: %#v", claims.RegisteredClaims)
	}
}

func TestJWTServiceRejectsExpiredWrongSignatureAndMalformedTokens(t *testing.T) {
	now := time.Date(2026, 6, 18, 4, 0, 0, 0, time.UTC)
	currentTime := now
	service := NewJWTService(
		"01234567890123456789012345678901",
		"api-integrator-gateway",
		time.Hour,
		func() time.Time { return currentTime },
	)
	token, _, err := service.Generate(model.User{
		ID:       7,
		Username: "marketplace",
		Role:     model.RoleAppUser,
		AppName:  "Marketplace",
	})
	if err != nil {
		t.Fatalf("Generate() returned an unexpected error: %v", err)
	}

	currentTime = now.Add(2 * time.Hour)
	if _, err := service.Validate(token); err == nil {
		t.Fatal("Validate() accepted expired token")
	}

	otherService := NewJWTService(
		"abcdefghijklmnopqrstuvwxyz123456",
		"api-integrator-gateway",
		time.Hour,
		func() time.Time { return now },
	)
	if _, err := otherService.Validate(token); err == nil {
		t.Fatal("Validate() accepted token with wrong signature")
	}
	if _, err := service.Validate("not-a-jwt"); err == nil {
		t.Fatal("Validate() accepted malformed token")
	}
}
