package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (repository *MySQLUserRepository) FindByUsernameAndApp(
	ctx context.Context,
	username string,
	appName string,
) (model.User, error) {
	row := repository.db.QueryRowContext(
		ctx,
		"SELECT id, username, password_hash, role, app_name FROM users WHERE username = ? AND app_name = ? LIMIT 1",
		username,
		appName,
	)
	return scanUser(row, "find user")
}

func (repository *MySQLUserRepository) FindFirstByRole(ctx context.Context, role model.Role) (model.User, error) {
	row := repository.db.QueryRowContext(
		ctx,
		"SELECT id, username, password_hash, role, app_name FROM users WHERE role = ? ORDER BY id ASC LIMIT 1",
		role,
	)
	return scanUser(row, "find first user by role")
}

func (repository *MySQLUserRepository) ListByRole(ctx context.Context, role model.Role) ([]model.User, error) {
	rows, err := repository.db.QueryContext(
		ctx,
		"SELECT id, username, password_hash, role, app_name FROM users WHERE role = ? ORDER BY app_name ASC, username ASC",
		role,
	)
	if err != nil {
		return nil, fmt.Errorf("list users by role: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		user, err := scanUser(rows, "scan user")
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (repository *MySQLUserRepository) Upsert(
	ctx context.Context,
	user model.User,
) error {
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO users (username, password_hash, role, app_name)
		 VALUES (?, ?, ?, ?)
		 ON DUPLICATE KEY UPDATE
		   password_hash = VALUES(password_hash),
		   role = VALUES(role)`,
		user.Username,
		user.PasswordHash,
		user.Role,
		user.AppName,
	)
	if err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}
	return nil
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(row userScanner, operation string) (model.User, error) {
	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.AppName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", operation, err)
	}
	return user, nil
}
