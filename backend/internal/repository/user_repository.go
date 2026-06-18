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
	var user model.User
	err := repository.db.QueryRowContext(
		ctx,
		"SELECT id, username, password_hash, role, app_name FROM users WHERE username = ? AND app_name = ? LIMIT 1",
		username,
		appName,
	).Scan(
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
		return model.User{}, fmt.Errorf("find user: %w", err)
	}
	return user, nil
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
