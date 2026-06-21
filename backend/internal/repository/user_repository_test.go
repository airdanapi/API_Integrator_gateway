package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

func TestMySQLUserRepositoryFindsUserByUsernameAndApp(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New(): %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, password_hash, role, app_name FROM users WHERE username = ? AND app_name = ? LIMIT 1",
	)).
		WithArgs("admin", "API Gateway").
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "username", "password_hash", "role", "app_name"},
		).AddRow(1, "admin", "hash", "admin_gateway", "API Gateway"))

	repo := NewMySQLUserRepository(db)
	user, err := repo.FindByUsernameAndApp(context.Background(), "admin", "API Gateway")
	if err != nil {
		t.Fatalf("FindByUsernameAndApp() returned an unexpected error: %v", err)
	}
	if user.ID != 1 ||
		user.Username != "admin" ||
		user.PasswordHash != "hash" ||
		user.Role != model.RoleAdminGateway ||
		user.AppName != "API Gateway" {
		t.Fatalf("unexpected user: %#v", user)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestMySQLUserRepositoryMapsMissingUserToDomainError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New(): %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, username, password_hash, role, app_name FROM users").
		WithArgs("missing", "Marketplace").
		WillReturnError(sql.ErrNoRows)

	repo := NewMySQLUserRepository(db)
	_, err = repo.FindByUsernameAndApp(
		context.Background(),
		"missing",
		"Marketplace",
	)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("FindByUsernameAndApp() error = %v, want ErrUserNotFound", err)
	}
}

func TestMySQLUserRepositoryUpsertsByUsernameAndApp(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New(): %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").
		WithArgs("admin", "hash", "admin_gateway", "API Gateway").
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewMySQLUserRepository(db)
	err = repo.Upsert(context.Background(), model.User{
		Username:     "admin",
		PasswordHash: "hash",
		Role:         model.RoleAdminGateway,
		AppName:      "API Gateway",
	})
	if err != nil {
		t.Fatalf("Upsert() returned an unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestMySQLUserRepositoryListByRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New(): %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, password_hash, role, app_name FROM users WHERE role = ? ORDER BY app_name ASC, username ASC",
	)).
		WithArgs(model.RoleAppUser).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "username", "password_hash", "role", "app_name"},
		).AddRow(2, "marketplace", "hash", "app_user", "Marketplace"))

	repo := NewMySQLUserRepository(db)
	users, err := repo.ListByRole(context.Background(), model.RoleAppUser)
	if err != nil {
		t.Fatalf("ListByRole() returned an unexpected error: %v", err)
	}
	if len(users) != 1 || users[0].Username != "marketplace" || users[0].AppName != "Marketplace" {
		t.Fatalf("unexpected users: %#v", users)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestMySQLUserRepositoryFindFirstByRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New(): %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, password_hash, role, app_name FROM users WHERE role = ? ORDER BY id ASC LIMIT 1",
	)).
		WithArgs(model.RoleAdminGateway).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "username", "password_hash", "role", "app_name"},
		).AddRow(1, "admin", "hash", "admin_gateway", "API Gateway"))

	repo := NewMySQLUserRepository(db)
	user, err := repo.FindFirstByRole(context.Background(), model.RoleAdminGateway)
	if err != nil {
		t.Fatalf("FindFirstByRole() returned an unexpected error: %v", err)
	}
	if user.Username != "admin" || user.Role != model.RoleAdminGateway {
		t.Fatalf("unexpected user: %#v", user)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}
