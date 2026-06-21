package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

func TestNotificationRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	n := model.Notification{
		CreatedAt: time.Now().UTC(),
		AppName:   "Marketplace",
		Type:      model.NotificationTypeAPIInactive,
		Message:   "API Marketplace tidak aktif selama lebih dari 1 minggu.",
		IsRead:    false,
	}

	mock.ExpectExec(`INSERT INTO notifications`).
		WithArgs(sqlmock.AnyArg(), n.AppName, n.Type, n.Message, 0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewMySQLNotificationRepository(db)
	id, err := repo.Insert(context.Background(), n)
	if err != nil {
		t.Fatalf("Insert() error: %v", err)
	}
	if id != 1 {
		t.Errorf("Insert() id = %d, want 1", id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}

func TestNotificationRepository_CountUnread(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM notifications WHERE app_name = \? AND is_read = 0`).
		WithArgs("Marketplace").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	repo := NewMySQLNotificationRepository(db)
	count, err := repo.CountUnread(context.Background(), "Marketplace")
	if err != nil {
		t.Fatalf("CountUnread() error: %v", err)
	}
	if count != 3 {
		t.Errorf("CountUnread() = %d, want 3", count)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}

func TestNotificationRepository_MarkAsRead(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`UPDATE notifications SET is_read = 1 WHERE id = \?`).
		WithArgs(int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := NewMySQLNotificationRepository(db)
	if err := repo.MarkAsRead(context.Background(), 5); err != nil {
		t.Fatalf("MarkAsRead() error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}

func TestNotificationRepository_ListAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	createdAt := time.Date(2026, 6, 21, 10, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`SELECT id, created_at, app_name, type, message, is_read\s+FROM notifications\s+ORDER BY created_at DESC LIMIT \? OFFSET \?`).
		WithArgs(10, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "app_name", "type", "message", "is_read"}).
			AddRow(1, createdAt, "Marketplace", model.NotificationTypeAPIInactive, "inactive", 0))

	repo := NewMySQLNotificationRepository(db)
	items, err := repo.ListAll(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("ListAll() error: %v", err)
	}
	if len(items) != 1 || items[0].AppName != "Marketplace" || items[0].IsRead {
		t.Fatalf("ListAll() items = %#v", items)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}

func TestNotificationRepository_CountUnreadAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM notifications WHERE is_read = 0`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(4))

	repo := NewMySQLNotificationRepository(db)
	count, err := repo.CountUnreadAll(context.Background())
	if err != nil {
		t.Fatalf("CountUnreadAll() error: %v", err)
	}
	if count != 4 {
		t.Errorf("CountUnreadAll() = %d, want 4", count)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}

func TestNotificationRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	createdAt := time.Date(2026, 6, 21, 10, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`SELECT id, created_at, app_name, type, message, is_read\s+FROM notifications WHERE id = \? LIMIT 1`).
		WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "app_name", "type", "message", "is_read"}).
			AddRow(7, createdAt, "POS", model.NotificationTypeErrorRate, "error rate", 1))

	repo := NewMySQLNotificationRepository(db)
	item, err := repo.FindByID(context.Background(), 7)
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if item.ID != 7 || item.AppName != "POS" || !item.IsRead {
		t.Fatalf("FindByID() = %#v", item)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}

func TestNotificationRepository_MarkAllAsReadAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`UPDATE notifications SET is_read = 1 WHERE is_read = 0`).
		WillReturnResult(sqlmock.NewResult(0, 3))

	repo := NewMySQLNotificationRepository(db)
	if err := repo.MarkAllAsReadAll(context.Background()); err != nil {
		t.Fatalf("MarkAllAsReadAll() error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}

func TestNotificationRepository_ExistsRecent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	since := time.Date(2026, 6, 20, 10, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM notifications WHERE app_name = \? AND type = \? AND created_at >= \?\)`).
		WithArgs("Marketplace", model.NotificationTypeAPIInactive, since).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	repo := NewMySQLNotificationRepository(db)
	exists, err := repo.ExistsRecent(context.Background(), "Marketplace", model.NotificationTypeAPIInactive, since)
	if err != nil {
		t.Fatalf("ExistsRecent() error: %v", err)
	}
	if !exists {
		t.Fatal("ExistsRecent() = false, want true")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}
