package auth

import (
	"context"
	"fmt"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

type SeedUser struct {
	Username string
	Password string
	Role     model.Role
	AppName  string
}

type UserWriter interface {
	Upsert(ctx context.Context, user model.User) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

func SeedUsers(
	ctx context.Context,
	users UserWriter,
	passwords PasswordHasher,
	seedUsers []SeedUser,
) error {
	for _, seed := range seedUsers {
		passwordHash, err := passwords.Hash(seed.Password)
		if err != nil {
			return fmt.Errorf("hash password for %s: %w", seed.Username, err)
		}
		if err := users.Upsert(ctx, model.User{
			Username:     seed.Username,
			PasswordHash: passwordHash,
			Role:         seed.Role,
			AppName:      seed.AppName,
		}); err != nil {
			return fmt.Errorf("seed user %s: %w", seed.Username, err)
		}
	}
	return nil
}
