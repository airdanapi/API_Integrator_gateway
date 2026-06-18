package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/repository"
)

type stubUserReader struct {
	user model.User
	err  error
}

func (stub stubUserReader) FindByUsernameAndApp(
	context.Context,
	string,
	string,
) (model.User, error) {
	return stub.user, stub.err
}

type stubPasswordVerifier struct {
	err error
}

func (stub stubPasswordVerifier) Compare(string, string) error {
	return stub.err
}

type stubTokenIssuer struct {
	token     string
	expiresIn int64
	err       error
}

func (stub stubTokenIssuer) Generate(model.User) (string, int64, error) {
	return stub.token, stub.expiresIn, stub.err
}

func TestServiceLoginReturnsRoleSpecificDashboard(t *testing.T) {
	tests := []struct {
		name          string
		role          model.Role
		wantDashboard string
	}{
		{name: "admin", role: model.RoleAdminGateway, wantDashboard: "/dashboard/admin"},
		{name: "app user", role: model.RoleAppUser, wantDashboard: "/dashboard/user"},
		{name: "monitoring", role: model.RoleMonitoringUser, wantDashboard: "/dashboard/monitoring"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewService(
				stubUserReader{user: model.User{
					ID:           1,
					Username:     "user",
					PasswordHash: "hash",
					Role:         test.role,
					AppName:      "Application",
				}},
				stubPasswordVerifier{},
				stubTokenIssuer{token: "signed-token", expiresIn: 3600},
			)

			result, err := service.Login(context.Background(), LoginRequest{
				Username: "user",
				Password: "password",
				AppName:  "Application",
			})
			if err != nil {
				t.Fatalf("Login() returned an unexpected error: %v", err)
			}
			if result.Token != "signed-token" ||
				result.Role != test.role ||
				result.AppName != "Application" ||
				result.DashboardURL != test.wantDashboard ||
				result.ExpiresIn != 3600 {
				t.Fatalf("unexpected Login() result: %#v", result)
			}
		})
	}
}

func TestServiceLoginUsesOneGenericCredentialError(t *testing.T) {
	tests := []struct {
		name       string
		reader     stubUserReader
		verifier   stubPasswordVerifier
		tokenError error
	}{
		{
			name:   "unknown user",
			reader: stubUserReader{err: repository.ErrUserNotFound},
		},
		{
			name: "wrong password",
			reader: stubUserReader{user: model.User{
				PasswordHash: "hash",
				Role:         model.RoleAppUser,
			}},
			verifier: stubPasswordVerifier{err: errors.New("mismatch")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewService(
				test.reader,
				test.verifier,
				stubTokenIssuer{err: test.tokenError},
			)
			_, err := service.Login(context.Background(), LoginRequest{
				Username: "user",
				Password: "password",
				AppName:  "Application",
			})
			if !errors.Is(err, ErrInvalidCredentials) {
				t.Fatalf("Login() error = %v, want ErrInvalidCredentials", err)
			}
		})
	}
}
