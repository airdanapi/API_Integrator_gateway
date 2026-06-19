package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
)

func TestLogRepository_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	durationMS := 100
	log := model.RequestLog{
		Timestamp:  time.Now().UTC(),
		SourceApp:  "Marketplace",
		Endpoint:   "/gateway/payment",
		Method:     "POST",
		Payload:    []byte(`{"amount":50000}`),
		Status:     200,
		Response:   []byte(`{"status":"success"}`),
		DurationMS: &durationMS,
	}

	mock.ExpectExec(`INSERT INTO request_logs`).
		WithArgs(
			sqlmock.AnyArg(), // timestamp
			log.SourceApp,
			log.Endpoint,
			log.Method,
			sqlmock.AnyArg(), // payload NullString
			log.Status,
			sqlmock.AnyArg(), // response NullString
			log.DurationMS,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := NewMySQLLogRepository(db)
	id, err := repo.Insert(context.Background(), log)
	if err != nil {
		t.Fatalf("Insert() returned unexpected error: %v", err)
	}
	if id != 1 {
		t.Errorf("Insert() id = %d, want 1", id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}

func TestLogRepository_ListRecent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	ts := time.Now().UTC()
	columns := []string{"id", "timestamp", "source_app", "endpoint", "method", "payload", "status", "response", "duration_ms"}
	rows := sqlmock.NewRows(columns).
		AddRow(1, ts, "Marketplace", "/gateway/payment", "POST", `{"amount":50000}`, 200, `{"status":"success"}`, 100).
		AddRow(2, ts, "POS", "/gateway/payment", "POST", nil, 500, nil, nil)

	mock.ExpectQuery(`SELECT id, timestamp, source_app, endpoint, method, payload, status, response, duration_ms FROM request_logs ORDER BY timestamp DESC`).
		WithArgs(10, 0).
		WillReturnRows(rows)

	repo := NewMySQLLogRepository(db)
	logs, err := repo.ListRecent(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("ListRecent() error: %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("ListRecent() count = %d, want 2", len(logs))
	}
	if logs[0].SourceApp != "Marketplace" {
		t.Errorf("logs[0].SourceApp = %q, want %q", logs[0].SourceApp, "Marketplace")
	}
	if logs[1].DurationMS != nil {
		t.Errorf("logs[1].DurationMS = %v, want nil", logs[1].DurationMS)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}

func TestLogRepository_CountByStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	since := time.Now().Add(-24 * time.Hour).UTC()
	rows := sqlmock.NewRows([]string{"status", "count"}).
		AddRow(200, 40).
		AddRow(400, 5).
		AddRow(500, 3)

	mock.ExpectQuery(`SELECT status, COUNT\(\*\) FROM request_logs WHERE timestamp >= \? GROUP BY status`).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	repo := NewMySQLLogRepository(db)
	counts, err := repo.CountByStatus(context.Background(), since)
	if err != nil {
		t.Fatalf("CountByStatus() error: %v", err)
	}
	if counts[200] != 40 {
		t.Errorf("counts[200] = %d, want 40", counts[200])
	}
	if counts[500] != 3 {
		t.Errorf("counts[500] = %d, want 3", counts[500])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}

func TestLogRepository_CountByStatusForApp(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	since := time.Now().Add(-7 * 24 * time.Hour).UTC()
	rows := sqlmock.NewRows([]string{"status", "count"}).
		AddRow(200, 12).
		AddRow(500, 3)

	mock.ExpectQuery(`SELECT status, COUNT\(\*\) FROM request_logs WHERE source_app = \? AND timestamp >= \? GROUP BY status`).
		WithArgs("Marketplace", sqlmock.AnyArg()).
		WillReturnRows(rows)

	repo := NewMySQLLogRepository(db)
	counts, err := repo.CountByStatusForApp(context.Background(), "Marketplace", since)
	if err != nil {
		t.Fatalf("CountByStatusForApp() error: %v", err)
	}
	if counts[200] != 12 {
		t.Errorf("counts[200] = %d, want 12", counts[200])
	}
	if counts[500] != 3 {
		t.Errorf("counts[500] = %d, want 3", counts[500])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}
