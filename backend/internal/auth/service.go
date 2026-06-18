package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	AppName  string `json:"app_name"`
}

type LoginResult struct {
	Token        string     `json:"token"`
	Role         model.Role `json:"role"`
	AppName      string     `json:"app_name"`
	DashboardURL string     `json:"dashboard_url"`
	ExpiresIn    int64      `json:"expires_in"`
}

type UserReader interface {
	FindByUsernameAndApp(
		ctx context.Context,
		username string,
		appName string,
	) (model.User, error)
}

type PasswordVerifier interface {
	Compare(passwordHash string, password string) error
}

type TokenIssuer interface {
	Generate(user model.User) (token string, expiresIn int64, err error)
}

type Service struct {
	users    UserReader
	password PasswordVerifier
	tokens   TokenIssuer
}

func NewService(
	users UserReader,
	password PasswordVerifier,
	tokens TokenIssuer,
) *Service {
	return &Service{
		users:    users,
		password: password,
		tokens:   tokens,
	}
}

func (service *Service) Login(
	ctx context.Context,
	request LoginRequest,
) (LoginResult, error) {
	user, err := service.users.FindByUsernameAndApp(
		ctx,
		request.Username,
		request.AppName,
	)
	if errors.Is(err, repository.ErrUserNotFound) {
		_ = service.password.Compare(dummyPasswordHash, request.Password)
		return LoginResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return LoginResult{}, fmt.Errorf("find login user: %w", err)
	}
	if err := service.password.Compare(user.PasswordHash, request.Password); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	dashboardURL, err := dashboardForRole(user.Role)
	if err != nil {
		return LoginResult{}, err
	}
	token, expiresIn, err := service.tokens.Generate(user)
	if err != nil {
		return LoginResult{}, fmt.Errorf("issue login token: %w", err)
	}
	return LoginResult{
		Token:        token,
		Role:         user.Role,
		AppName:      user.AppName,
		DashboardURL: dashboardURL,
		ExpiresIn:    expiresIn,
	}, nil
}

func dashboardForRole(role model.Role) (string, error) {
	switch role {
	case model.RoleAdminGateway:
		return "/dashboard/admin", nil
	case model.RoleAppUser:
		return "/dashboard/user", nil
	case model.RoleMonitoringUser:
		return "/dashboard/monitoring", nil
	default:
		return "", fmt.Errorf("unsupported role %q", role)
	}
}
