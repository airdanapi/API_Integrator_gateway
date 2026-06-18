package database

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMigrateCreatesUsersSchemaAndIsIdempotent(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN is not set")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("sql.Open(): %v", err)
	}
	defer db.Close()

	if err := Migrate(context.Background(), db); err != nil {
		t.Fatalf("first Migrate() returned an unexpected error: %v", err)
	}
	if err := Migrate(context.Background(), db); err != nil {
		t.Fatalf("second Migrate() returned an unexpected error: %v", err)
	}

	var databaseName string
	if err := db.QueryRow("SELECT DATABASE()").Scan(&databaseName); err != nil {
		t.Fatalf("read current database: %v", err)
	}

	var columnCount int
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_schema = ?
		  AND table_name = 'users'
		  AND column_name IN ('id', 'username', 'password_hash', 'role', 'app_name')
	`, databaseName).Scan(&columnCount)
	if err != nil {
		t.Fatalf("inspect users columns: %v", err)
	}
	if columnCount != 5 {
		t.Fatalf("users required column count = %d, want 5", columnCount)
	}

	var uniqueColumnCount int
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.statistics
		WHERE table_schema = ?
		  AND table_name = 'users'
		  AND index_name = 'uq_users_username_app_name'
		  AND non_unique = 0
	`, databaseName).Scan(&uniqueColumnCount)
	if err != nil {
		t.Fatalf("inspect users unique index: %v", err)
	}
	if uniqueColumnCount != 2 {
		t.Fatalf("unique index column count = %d, want 2", uniqueColumnCount)
	}
}
