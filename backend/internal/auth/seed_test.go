package auth

import (
	"context"
	"testing"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

type memoryUserWriter struct {
	users map[string]model.User
}

func (writer *memoryUserWriter) Upsert(_ context.Context, user model.User) error {
	if writer.users == nil {
		writer.users = make(map[string]model.User)
	}
	writer.users[user.Username+"|"+user.AppName] = user
	return nil
}

func TestSeedUsersIsIdempotentAndHashesEveryPassword(t *testing.T) {
	writer := &memoryUserWriter{}
	hasher := NewBcryptPasswordHasher()
	users := []SeedUser{
		{
			Username: "admin",
			Password: "admin-password",
			Role:     model.RoleAdminGateway,
			AppName:  "API Gateway",
		},
		{
			Username: "marketplace",
			Password: "marketplace-password",
			Role:     model.RoleAppUser,
			AppName:  "Marketplace",
		},
		{
			Username: "insight",
			Password: "insight-password",
			Role:     model.RoleMonitoringUser,
			AppName:  "UMKM Insight",
		},
	}

	if err := SeedUsers(context.Background(), writer, hasher, users); err != nil {
		t.Fatalf("first SeedUsers() returned an unexpected error: %v", err)
	}
	if err := SeedUsers(context.Background(), writer, hasher, users); err != nil {
		t.Fatalf("second SeedUsers() returned an unexpected error: %v", err)
	}
	if len(writer.users) != 3 {
		t.Fatalf("seeded user count = %d, want 3", len(writer.users))
	}

	for _, seed := range users {
		stored := writer.users[seed.Username+"|"+seed.AppName]
		if stored.PasswordHash == "" || stored.PasswordHash == seed.Password {
			t.Fatalf("password for %s was not hashed", seed.Username)
		}
		if err := hasher.Compare(stored.PasswordHash, seed.Password); err != nil {
			t.Fatalf("stored password for %s does not match seed: %v", seed.Username, err)
		}
	}
}
