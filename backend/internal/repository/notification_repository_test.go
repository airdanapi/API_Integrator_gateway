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
